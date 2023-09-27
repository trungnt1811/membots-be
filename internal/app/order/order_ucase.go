package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
	"gorm.io/datatypes"
)

type orderUCase struct {
	Repo     interfaces.OrderRepository
	ATRepo   interfaces.ATRepository
	Producer *msgqueue.QueueWriter
}

func NewOrderUCase(repo interfaces.OrderRepository, atRepo interfaces.ATRepository) interfaces.OrderUCase {
	return &orderUCase{
		Repo:     repo,
		ATRepo:   atRepo,
		Producer: msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_AFF_ORDER_UPDATE),
	}
}

func (u *orderUCase) PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error) {
	// First, save log
	bytes, err := json.Marshal(postBackReq)
	if err != nil {
		log.LG.Errorf("cannot marshall post-back req: %v", err)
	} else {
		err := u.Repo.SavePostBackLog(&model.AffPostBackLog{
			OrderId:   postBackReq.OrderId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Data:      datatypes.JSON(bytes),
		})
		if err != nil {
			log.LG.Errorf("save aff_post_back_log error: %v", err)
		}
	}

	// Get campaign to make sure supported
	campaign, err := u.Repo.GetCampaignByATId(postBackReq.CampaignId)
	if err != nil {
		return nil, fmt.Errorf("post_back campaign: %v", err)
	}

	// Parse sale time for later request
	salesTime, err := util.ParsePostBackTime(postBackReq.SalesTime)
	if err != nil {
		return nil, fmt.Errorf("cannot parse sales_time: %v", err)
	}
	since, until := util.GetSinceUntilTime(salesTime, 1)

	// Then, find if order is exist
	resp, err := u.ATRepo.QueryOrders(types.ATOrderQuery{
		Since: since,
		Until: until,
	}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("query order-list error: %v", err)
	}
	var atOrder *types.ATOrder
	for idx, item := range resp.Data {
		if item.OrderId == postBackReq.OrderId {
			atOrder = &resp.Data[idx]
		}
	}
	if atOrder == nil {
		return nil, fmt.Errorf("not found order \"%s\" in order-list", postBackReq.OrderId)
	}

	// Parse user_id, tracked_id from utm_content
	userId, trackedId := util.ParseUTMContent(atOrder.UTMContent)

	order, err := u.Repo.FindOrderByAccessTradeId(atOrder.OrderId)
	statusChanged := false
	if err != nil {
		if err.Error() == "record not found" {
			// Order not exist, create new one
			order = model.NewOrderFromATOrder(userId, campaign.ID, campaign.BrandId, atOrder)
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()
			crErr := u.Repo.CreateOrder(order)
			if crErr != nil {
				return nil, fmt.Errorf("create order failed: %v", crErr)
			}
			// When new order created, mark as status changed
			statusChanged = true
		} else {
			return nil, err
		}
	} else {
		// Or update exist order
		updated := model.NewOrderFromATOrder(userId, campaign.ID, campaign.BrandId, atOrder)
		updated.UpdatedAt = time.Now()
		updated.ID = order.ID

		// When order is updated, check if status changed or not
		statusChanged = order.CheckStatusChanged(updated)

		_, err = u.Repo.UpdateOrder(updated)
		if err != nil {
			return nil, fmt.Errorf("update order failed: %v", err)
		}
		order = updated
	}

	// After create order, push Kafka msg for sync transactions
	_, err = u.SyncTransactionsByOrder(order.AccessTradeOrderId)
	if err != nil {
		log.LG.Errorf("sync txs failed: %v", err)
	}

	// Send Kafka message if order status changed
	if statusChanged {
		// Order has been approved
		u.sendOrderUpdateMsg(atOrder.OrderId, order.OrderStatus, atOrder.IsConfirmed)
	}

	// And update tracked item
	err = u.Repo.UpdateTrackedClickOrder(trackedId, order)
	if err != nil {
		log.LG.Errorf("update tracked click failed: %v", err)
	}

	return order, nil
}

// SyncTransactionsByOrder The `SyncTransactionsByOrder` function is responsible for synchronizing transactions for a given
// order.
// TODO: Sync txs every day instead of order time
func (u *orderUCase) SyncTransactionsByOrder(atOrderId string) (int, error) {
	// First find created order
	order, err := u.Repo.FindOrderByAccessTradeId(atOrderId)
	if err != nil {
		return 0, fmt.Errorf("find order id \"%s\" failed: %v", atOrderId, err)
	}
	since, until := util.GetSinceUntilTime(order.SalesTime, 2)
	txs, err := u.ATRepo.QueryTransactions(types.ATTransactionQuery{
		Since:         since,
		Until:         until,
		TransactionId: atOrderId, // NOTE: post back req will return old transaction_id as order_id!
	}, 0, 0)
	if err != nil {
		return 0, fmt.Errorf("query txs failed: %v", err)
	}
	if len(txs.Data) == 0 {
		return 0, fmt.Errorf("txs empty for order: %s", atOrderId)
	}

	// Save transactions if found
	transactions := make([]model.AffTransaction, len(txs.Data))
	for i, it := range txs.Data {
		transactions[i] = *model.NewAffTransactionFromAT(order, &it)
	}

	err = u.Repo.UpdateOrCreateATTransactions(transactions)
	if err != nil {
		return 0, fmt.Errorf("save txs failed: %v", err)
	}

	return len(txs.Data), nil
}

func (u *orderUCase) GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*dto.OrderDetailsDto, error) {
	order, err := u.Repo.GetOrderDetails(ctx, userId, orderId)
	if err != nil {
		return nil, err
	}

	orderDto := order.ToOrderDetailsDto()
	return &orderDto, nil
}

func (u *orderUCase) GetOrderHistory(ctx context.Context, userId uint32, status string, page, size int) (dto.OrderHistoryResponse, error) {
	pastTimeLimit := time.Now().Add(-6 * 30 * 24 * time.Hour) // 6 months before - user cannot query order older than this time
	orderHistory, err := u.Repo.GetOrderHistory(ctx, pastTimeLimit, userId, status, page, size)
	if err != nil {
		return dto.OrderHistoryResponse{}, err
	}

	nextPage := page
	if len(orderHistory) > size {
		nextPage = page + 1
	}

	orderDtos := []dto.OrderDetailsDto{}
	for i, item := range orderHistory {
		if i >= size {
			break
		}
		orderDtos = append(orderDtos, item.ToOrderDetailsDto())
	}

	totalOrder, err := u.Repo.CountOrders(ctx, pastTimeLimit, userId, status)
	if err != nil {
		return dto.OrderHistoryResponse{}, err
	}

	return dto.OrderHistoryResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     orderDtos,
		Total:    totalOrder,
	}, nil
}

func (u *orderUCase) CheckOrderConfirmed() (int, error) {
	// First query approved order which is not confirmed
	t := time.Now()
	q := map[string]any{
		"is_confirmed": 0,
		"order_status": model.OrderStatusApproved,
	}
	orders, err := u.Repo.QueryOrdersConfirmedBefore(t, q)
	if err != nil {
		return 0, err
	}

	if len(orders) == 0 {
		return 0, nil
	}

	updatedCount := 0

	mappedATOrders := map[string]types.ATOrder{}
	for _, order := range orders {
		atOrder, ok := mappedATOrders[order.AccessTradeOrderId]
		if !ok {
			// Fetch from AT
			since, until := util.GetSinceUntilTime(order.SalesTime, 2)
			resp, err := u.ATRepo.QueryOrders(types.ATOrderQuery{
				Since: since,
				Until: until,
			}, 0, 0)
			if err != nil {
				log.LG.Errorf("query orders from accesstrade failed: %v", err)
				continue
			}

			for idx, item := range resp.Data {
				mappedATOrders[item.OrderId] = resp.Data[idx]
				if item.OrderId == order.AccessTradeOrderId {
					atOrder = resp.Data[idx]
				}
			}
		}

		if atOrder.OrderId == "" {
			// Still empty
			log.LG.Errorf("cannot found access trade order: %s", order.AccessTradeOrderId)
		}

		// Check if confirmed or not
		if atOrder.IsConfirmed == 0 {
			// Update order as cancelled
			order.OrderStatus = model.OrderStatusCancelled

			_, err := u.Repo.UpdateOrder(&order)
			if err != nil {
				log.LG.Errorf("update order cancelled err: %v", err)
			}
			// Send msg to Kafka
			u.sendOrderUpdateMsg(atOrder.OrderId, order.OrderStatus, atOrder.IsConfirmed)
		}
	}

	return updatedCount, nil
}

func (u *orderUCase) sendOrderUpdateMsg(orderId string, orderStatus string, isConfirmed uint8) {
	// Order has been approved
	msg := msgqueue.MsgOrderUpdated{
		AtOrderID:   orderId,
		OrderStatus: orderStatus,
		IsConfirmed: isConfirmed,
	}
	msgValue, err := json.Marshal(&msg)
	if err != nil {
		log.LG.Errorf("marshall order approved error: %v", err)
	} else {
		err = u.Producer.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte(orderId),
				Value: msgValue,
			},
		)
		if err != nil {
			log.LG.Errorf("produce order approved msg failed: %v", err)
		}
	}
}

package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
	"gorm.io/datatypes"
)

type OrderUcase struct {
	Repo     interfaces.OrderRepository
	ATRepo   interfaces.ATRepository
	Producer *msgqueue.QueueWriter
}

func NewOrderUcase(repo interfaces.OrderRepository, atRepo interfaces.ATRepository) *OrderUcase {
	return &OrderUcase{
		Repo:     repo,
		ATRepo:   atRepo,
		Producer: msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_AFF_ORDER_APPROVE),
	}
}

func (u *OrderUcase) PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error) {
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
	if err != nil {
		if err.Error() == "record not found" {
			// Order not exist, create new one
			order = model.NewOrderFromATOrder(userId, atOrder)
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()
			crErr := u.Repo.CreateOrder(order)
			if crErr != nil {
				return nil, fmt.Errorf("create order failed: %v", crErr)
			}
		} else {
			return nil, err
		}
	} else {
		// Or update exist order
		updated := model.NewOrderFromATOrder(userId, atOrder)
		order.UpdatedAt = time.Now()
		updated.ID = order.ID
		_, err = u.Repo.UpdateOrder(updated)
		if err != nil {
			return nil, fmt.Errorf("update order failed: %v", err)
		}
	}

	// After create order, push Kafka msg for sync transactions
	_, err = u.SyncTransactionsByOrder(order.AccessTradeOrderId)
	if err != nil {
		log.LG.Errorf("sync txs failed: %v", err)
	}
	// Send Kafka message if order approved
	if order.OrderApproved != 0 {
		// Order has been approved
		msgValue := []byte(fmt.Sprintf("{\"accesstrade_order_id\":%s}", atOrder.OrderId))
		err = u.Producer.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte(atOrder.OrderId),
				Value: msgValue,
			},
		)
		if err != nil {
			log.LG.Errorf("produce order approved msg failed: %v", err)
		}
	}

	// And update tracked item
	err = u.Repo.UpdateTrackedClickOrder(trackedId, order)
	if err != nil {
		log.LG.Errorf("update tracked click failed: %v", err)
	}

	return order, nil
}

// The `SyncTransactionsByOrder` function is responsible for synchronizing transactions for a given
// order.
// TODO: Sync txs every day instead of order time
func (u *OrderUcase) SyncTransactionsByOrder(atOrderId string) (int, error) {
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

func (u *OrderUcase) GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*dto.OrderDetailsDto, error) {
	order, err := u.Repo.GetOrderDetails(ctx, userId, orderId)
	if err != nil {
		return nil, err
	}
	orderDto := order.ToOrderDetailsDto()
	return &orderDto, nil
}

func (u *OrderUcase) GetOrderHistory(ctx context.Context, userId uint32, status string, page, size int) (dto.OrderHistoryResponse, error) {
	pastTimeLimit := time.Now().Add(-6 * 30 * 24 * time.Hour) // 6 months before - user cannot query order older than this time
	orderHistory, err := u.Repo.GetOrderHistory(ctx, pastTimeLimit, userId, status, page, size)
	if err != nil {
		return dto.OrderHistoryResponse{}, err
	}

	nextPage := page
	if len(orderHistory) > size {
		nextPage = page + 1
	}

	orderDtos := make([]dto.OrderDetailsDto, len(orderHistory))
	for idx, item := range orderHistory {
		orderDtos[idx] = item.ToOrderDetailsDto()
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

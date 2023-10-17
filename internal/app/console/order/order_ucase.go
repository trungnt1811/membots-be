package consoleOrder

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
	"golang.org/x/exp/slices"
)

type ConsoleOrderUcase struct {
	Repo                interfaces.ConsoleOrderRepository
	OrderRepo           interfaces.OrderRepository
	OrderUpdateProducer *msgqueue.QueueWriter
}

func NewConsoleOrderUcase(repo interfaces.ConsoleOrderRepository,
	orderRepo interfaces.OrderRepository,
	orderUpdateProducer *msgqueue.QueueWriter) *ConsoleOrderUcase {
	return &ConsoleOrderUcase{
		Repo:                repo,
		OrderRepo:           orderRepo,
		OrderUpdateProducer: orderUpdateProducer,
	}
}

func (u *ConsoleOrderUcase) GetOrderList(q *dto.OrderListQuery) (*dto.OrderListResponse, error) {
	dbQuery := map[string]any{}
	if q.OrderStatus != "" {
		dbQuery["order_status"] = q.OrderStatus
	}
	if q.UserId != 0 {
		dbQuery["user_id"] = q.UserId
	}

	timeRange := dto.TimeRange{
		Since: q.Since,
		Until: q.Until,
	}

	list, total, err := u.Repo.FindOrdersByQuery(timeRange, dbQuery, q.Page, q.PerPage)
	if err != nil {
		return nil, fmt.Errorf("find list failed: %v", err)
	}
	totalPages := 1.0
	if q.PerPage != 0 {
		totalPages = math.Ceil(float64(total) / float64(q.PerPage))
	}
	resp := dto.OrderListResponse{
		Total:      int(total),
		Page:       q.Page,
		PerPage:    q.PerPage,
		TotalPages: int(totalPages),
		Data:       make([]dto.AffOrder, len(list)),
	}
	for idx, item := range list {
		resp.Data[idx] = item.ToDto()
	}
	return &resp, err
}

func (u *ConsoleOrderUcase) GetOrderByOrderId(orderId string) (*dto.AffOrder, error) {
	order, txs, err := u.Repo.FindOrderByOrderId(orderId)
	if err != nil {
		return nil, fmt.Errorf("find order failed: %v", err)
	}

	affOrder := order.ToDto()
	affOrder.Transactions = make([]dto.AffTransaction, len(txs))
	for i, tx := range txs {
		affOrder.Transactions[i] = tx.ToDto()
	}

	return &affOrder, nil
}

func (u *ConsoleOrderUcase) SyncOrderReward(atOrderId string) error {
	order, err := u.OrderRepo.FindOrderByAccessTradeId(atOrderId)
	if err != nil {
		return err
	}

	validStatuses := []string{model.OrderStatusInitial, model.OrderStatusPending, model.OrderStatusApproved}
	if !slices.Contains(validStatuses, order.OrderStatus) {
		return fmt.Errorf("cannot sync reward amount if order not in status %v", validStatuses)
	}

	u.sendOrderUpdateMsg(order.UserId, order.AccessTradeOrderId, order.OrderStatus, order.IsConfirmed)

	return nil
}

func (u *ConsoleOrderUcase) sendOrderUpdateMsg(userId uint, orderId string, orderStatus string, isConfirmed uint8) {
	// Order has been approved
	msg := msgqueue.MsgOrderUpdated{
		UserId:      userId,
		AtOrderID:   orderId,
		OrderStatus: orderStatus,
		IsConfirmed: isConfirmed,
	}
	msgValue, err := json.Marshal(&msg)
	if err != nil {
		log.LG.Errorf("marshall order approved error: %v", err)
	} else {
		if u.OrderUpdateProducer == nil {
			log.LG.Error("produce is nil")
		} else {
			err = u.OrderUpdateProducer.WriteMessages(
				context.Background(),
				kafka.Message{
					Key:   []byte(orderId),
					Value: msgValue,
				},
			)
			if err != nil {
				log.LG.Errorf("sync order reward failed: %v", err)
			}
		}
	}
}

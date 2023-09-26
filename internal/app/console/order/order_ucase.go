package consoleOrder

import (
	"fmt"
	"math"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type ConsoleOrderUcase struct {
	Repo interfaces.ConsoleOrderRepository
}

func NewConsoleOrderUcase(repo interfaces.ConsoleOrderRepository) *ConsoleOrderUcase {
	return &ConsoleOrderUcase{
		Repo: repo,
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

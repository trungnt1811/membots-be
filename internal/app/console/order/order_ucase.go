package consoleOrder

import (
	"fmt"
	"math"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type OrderUcase struct {
	Repo interfaces.ConsoleOrderRepository
}

func NewOrderUcase(repo interfaces.ConsoleOrderRepository) *OrderUcase {
	return &OrderUcase{
		Repo: repo,
	}
}

func (u *OrderUcase) GetOrderList(q *dto.OrderListQuery) (*dto.OrderListResponse, error) {
	dbQuery := map[string]any{}
	if q.OrderStatus != "" {
		dbQuery["order_status"] = q.OrderStatus
	}
	if q.UserId != 0 {
		dbQuery["user_id"] = q.UserId
	}

	timeRange := dto.TimeRange{}
	if !q.Since.IsZero() {
		timeRange.Since = &q.Since
	}
	if !q.Until.IsZero() {
		timeRange.Until = &q.Until
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

package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type ConsoleOrderRepository interface {
	FindOrdersByQuery(timeRange dto.TimeRange, dbQuery map[string]any, page int, perPage int) ([]model.AffOrder, int64, error)
	FindOrderByOrderId(orderId string) (*model.AffOrder, []model.AffTransaction, error)
}

type ConsoleOrderUcase interface {
	GetOrderList(q *dto.OrderListQuery) (*dto.OrderListResponse, error)
	GetOrderByOrderId(orderId string) (*dto.AffOrder, error)
}

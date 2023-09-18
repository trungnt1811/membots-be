package interfaces

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/order/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type OrderRepository interface {
	FindNonRewardOrders(
		customerId, sellerId int,
		fromDate time.Time,
		minValue int64,
		additionalFilter map[string]interface{},
	) ([]model.AffOrder, error)

	SavePostBackLog(req *model.AffPostBackLog) error
	CreateOrder(order *model.AffOrder) error
	UpdateOrder(updated *model.AffOrder) (int, error)
	FindOrderByAccessTradeId(atOrderId string) (*model.AffOrder, error)
	UpdateOrCreateATTransactions([]model.AffTransaction) error

	UpdateTrackedClickOrder(trackedId uint64, order *model.AffOrder) error
	GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*dto.OrderDetailsDto, error)
	GetOrderHistory(ctx context.Context, userId uint32, page, size int) ([]dto.OrderDetailsDto, error)
	CountOrder(ctx context.Context, userId uint32) (int64, error)
	FindOrdersByQuery(timeRange types.TimeRange, dbQuery map[string]any, page int, perPage int) ([]model.AffOrder, int64, error)
}

type OrderUcase interface {
	PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error)
	GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*dto.OrderDetailsDto, error)
	GetOrderHistory(ctx context.Context, userId uint32, page, size int) (dto.OrderHistoryResponse, error)

	GetOrderList(q *dto.OrderListQuery) (*dto.OrderListResponse, error)
}

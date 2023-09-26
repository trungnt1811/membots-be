package interfaces

import (
	"context"
	"time"

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
	GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*model.OrderDetails, error)
	GetOrderHistory(ctx context.Context, since time.Time, userId uint32, status string, page, size int) ([]model.OrderDetails, error)
	CountOrders(ctx context.Context, since time.Time, userId uint32, status string) (int64, error)

	GetCampaignByATId(atId string) (*model.AffCampaign, error)
}

type OrderUcase interface {
	PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error)
	GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*dto.OrderDetailsDto, error)
	GetOrderHistory(ctx context.Context, userId uint32, status string, page, size int) (dto.OrderHistoryResponse, error)
}

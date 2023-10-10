package interfaces

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type OrderRepository interface {
	QueryOrdersConfirmedBefore(t time.Time, q map[string]any) ([]model.AffOrder, error)

	SavePostBackLog(req *model.AffPostBackLog) error
	CreateOrder(order *model.AffOrder) error
	UpdateOrder(updated *model.AffOrder) (int, error)
	FindOrderByAccessTradeId(atOrderId string) (*model.AffOrder, error)
	FindOrderMappedByAccessTradeIds(atOrderIds []string) (map[string]model.AffOrder, error)
	UpdateOrCreateATTransactions([]model.AffTransaction) error

	GetCampaignByTrackedClick(trackedId uint64) (*model.AffCampaign, error)
	UpdateTrackedClickOrder(trackedId uint64, order *model.AffOrder) error
	GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*model.OrderDetails, error)
	GetOrderHistory(ctx context.Context, since time.Time, userId uint32, status string, page, size int) ([]model.OrderDetails, error)
	CountOrders(ctx context.Context, since time.Time, userId uint32, status string) (int64, error)

	GetCampaignByATId(atId string) (*model.AffCampaign, error)
}

type OrderUCase interface {
	PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error)
	GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*dto.OrderDetailsDto, error)
	GetOrderHistory(ctx context.Context, userId uint32, status string, page, size int) (dto.OrderHistoryResponse, error)
	SyncTransactionsByOrder(atOrderId string) (int, error)

	CheckOrderConfirmed() (int, error)
	CheckOrderListAndSync() (int, error)
}

package interfaces

import (
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
}

type OrderUcase interface {
	PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error)
}

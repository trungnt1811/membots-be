package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/model"
	"time"
)

type OrderRepository interface {
	FindNonRewardOrders(
		customerId, sellerId int,
		fromDate time.Time,
		minValue int64,
		additionalFilter map[string]interface{},
	) ([]model.OrderEntity, error)
}

type OrderUCase interface{}

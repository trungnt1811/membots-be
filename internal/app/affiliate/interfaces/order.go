package interfaces

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
)

type OrderRepository interface {
	FindNonRewardOrders(
		customerId, sellerId int,
		fromDate time.Time,
		minValue int64,
		additionalFilter map[string]interface{},
	) ([]model.OrderEntity, error)
}

type OrderUsecase interface{}

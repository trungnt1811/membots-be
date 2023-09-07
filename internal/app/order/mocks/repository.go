package mocks

import (
	"github.com/astraprotocol/affiliate-system/internal/model"
	"time"
)

type MockOrderRepository struct {
	orders []model.OrderEntity
}

func NewMockOrderRepository(orders []model.OrderEntity) *MockOrderRepository {
	return &MockOrderRepository{
		orders: orders,
	}
}

func (mr *MockOrderRepository) FindNonRewardOrders(
	customerId, sellerId int,
	fromDate time.Time,
	minValue int64,
	additionalFilter map[string]interface{},
) ([]model.OrderEntity, error) {
	return mr.orders, nil
}

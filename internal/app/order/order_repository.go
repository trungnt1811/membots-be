package order

import (
	"github.com/astraprotocol/affiliate-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) FindNonRewardOrders(
	customerId, sellerId int,
	fromDate time.Time,
	minValue int64,
	additionalFilter map[string]interface{},
) ([]model.OrderEntity, error) {
	var entities []model.OrderEntity
	err := repo.db.
		Table("order").
		Where("customer_id = ?", customerId).
		Where("seller_id = ?", sellerId).
		Where("amount >= ?", minValue).
		Where("reward_id", nil).
		Where("initialized_at >= ?", fromDate).
		Find(&entities, additionalFilter).Error

	if err != nil {
		return nil, err
	}
	return entities, nil
}

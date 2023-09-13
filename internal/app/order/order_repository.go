package order

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/model"

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
) ([]model.AffOrder, error) {
	var entities []model.AffOrder
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

func (repo *OrderRepository) SavePostBackLog(req *model.AffPostBackLog) error {
	return repo.db.Create(req).Error
}

func (repo *OrderRepository) CreateOrder(order *model.AffOrder) error {
	return repo.db.Create(order).Error
}

func (repo *OrderRepository) UpdateOrder(updated *model.AffOrder) (int, error) {
	result := repo.db.Model(updated).Where("id = ?", updated.ID).Updates(updated)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (repo *OrderRepository) UpdateOrCreateATTransactions(newTxs []model.AffTransaction) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		for _, newTx := range newTxs {
			// Find by id
			var oldTx model.AffTransaction
			err := tx.First(&oldTx, "accesstrade_id = ?", newTx.AccessTradeId).Error
			if err != nil {
				if err.Error() == "record not found" {
					// If tx is not found, create one
					newTx.CreatedAt = time.Now()
					newTx.UpdatedAt = time.Now()
					crErr := tx.Create(&newTx).Error
					if crErr != nil {
						return crErr
					}
					continue
				} else {
					return err
				}
			}
			// When found one, compare and update
			newTx.ID = oldTx.ID
			newTx.UpdatedAt = time.Now()
			upErr := tx.Model(&oldTx).Updates(&newTx).Error
			if upErr != nil {
				return upErr
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *OrderRepository) FindOrderByAccessTradeId(atOrderId string) (*model.AffOrder, error) {
	var order model.AffOrder
	err := repo.db.First(&order, "accesstrade_order_id = ?", atOrderId).Error
	return &order, err
}

func (repo *OrderRepository) UpdateTrackedClickOrder(trackedId uint64, order *model.AffOrder) error {
	// Only update empty order_id item
	err := repo.db.Model(&model.AffTrackedClick{}).Where(map[string]any{
		"id":       trackedId,
		"order_id": "",
	}).Updates(map[string]any{
		"order_id":   order.AccessTradeOrderId,
		"aff_link":   order.AffLink,
		"updated_at": time.Now(),
	}).Error

	return err
}

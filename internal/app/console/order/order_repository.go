package consoleOrder

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type ConsoleOrderRepository struct {
	db *gorm.DB
}

func NewConsoleOrderRepository(db *gorm.DB) *ConsoleOrderRepository {
	return &ConsoleOrderRepository{
		db: db,
	}
}

func (repo *ConsoleOrderRepository) FindOrdersByQuery(timeRange dto.TimeRange, dbQuery map[string]any, page int, perPage int) ([]model.AffOrder, int64, error) {
	var orders []model.AffOrder
	tx := repo.db.Model(&orders)
	totalTx := repo.db.Model(&orders)
	if perPage != 0 {
		tx.Limit(perPage)

		if page != 0 {
			tx.Offset((page - 1) * perPage)
		}
	}

	if timeRange.Since != nil {
		tx.Where(
			"created_at >= ?", timeRange.Since,
		)
		totalTx.Where(
			"created_at >= ?", timeRange.Since,
		)
	}
	if timeRange.Until != nil {
		tx.Where(
			"created_at <= ?", timeRange.Until,
		)
		totalTx.Where(
			"created_at >= ?", timeRange.Since,
		)
	}

	tx.Where(dbQuery)
	totalTx.Where(dbQuery)

	err := tx.Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = totalTx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, err
}

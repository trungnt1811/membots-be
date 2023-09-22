package statistic

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type StatisticRepository struct {
	Db *gorm.DB
}

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{
		Db: db,
	}
}

func (repo *StatisticRepository) FindOrdersInRange(d dto.TimeRange, offset int, limit int) ([]model.AffOrder, error) {
	var orders []model.AffOrder
	m := repo.Db.Model(&orders)
	if d.Since != nil {
		m.Where("created_at >= ?", d.Since)
	}
	if d.Until != nil {
		m.Where("created_at <= ?", d.Until)
	}
	if offset != 0 {
		m.Offset(offset)
	}
	if limit != 0 {
		m.Limit(limit)
	}

	err := m.Find(&orders).Error
	return orders, err
}

func (repo *StatisticRepository) TotalOrdersInRange(d dto.TimeRange) (int64, error) {
	m := repo.Db.Model(&model.AffOrder{})
	if d.Since != nil {
		m.Where("created_at >= ?", d.Since)
	}
	if d.Until != nil {
		m.Where("created_at <= ?", d.Until)
	}

	var result int64 = 0
	err := m.Count(&result).Error
	return result, err
}

func (repo *StatisticRepository) CalculateCashbackInRange(d dto.TimeRange) (dto.Cashback, error) {
	var cashback dto.Cashback

	sql := repo.Db.Table("aff_reward").Select(
		"SUM(rewarded_amount) as distributed, SUM(amount) as remain",
	)
	if d.Since != nil {
		sql.Where("updated_at >= ?", d.Since)
	}
	if d.Until != nil {
		sql.Where("updated_at <= ?", d.Until)
	}

	err := sql.Row().Scan(&cashback.Distributed, &cashback.Remain)

	if err != nil {
		return cashback, err
	}
	cashback.Remain = cashback.Remain - cashback.Distributed
	return cashback, nil
}

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

func (repo *StatisticRepository) FindOrdersInRange(campaignId uint, d dto.TimeRange, offset int, limit int) ([]model.AffOrder, error) {
	var orders []model.AffOrder
	m := repo.Db.Model(&orders)
	if !d.Since.IsZero() {
		m.Where("created_at >= ?", d.Since)
	}
	if !d.Until.IsZero() {
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

func (repo *StatisticRepository) TotalOrdersInRange(campaignId uint, d dto.TimeRange) (int64, error) {
	m := repo.Db.Model(&model.AffOrder{})
	if !d.Since.IsZero() {
		m.Where("created_at >= ?", d.Since)
	}
	if !d.Until.IsZero() {
		m.Where("created_at <= ?", d.Until)
	}
	if campaignId != 0 {
		m.Where("campaign_id = ?", campaignId)
	}

	var result int64 = 0
	err := m.Count(&result).Error
	return result, err
}

func (repo *StatisticRepository) TotalActiveCampaignsInRange(d dto.TimeRange) (int64, error) {
	m := repo.Db.Model(&model.AffOrder{})
	if !d.Since.IsZero() {
		m.Where("created_at >= ?", d.Since)
	}
	if !d.Until.IsZero() {
		m.Where("created_at <= ?", d.Until)
	}

	var result int64 = 0
	err := m.Distinct("campaign_id").Count(&result).Error
	return result, err
}

func (repo *StatisticRepository) CalculateCashbackInRange(campaignId uint, d dto.TimeRange) (dto.Cashback, error) {
	var cashback dto.Cashback

	sql := repo.Db.Table("aff_reward").Select(
		"SUM(rewarded_amount) as distributed, SUM(amount) as remain",
	)
	if !d.Since.IsZero() {
		sql.Where("updated_at >= ?", d.Since)
	}
	if !d.Until.IsZero() {
		sql.Where("updated_at <= ?", d.Until)
	}

	if campaignId != 0 {
		sql.Joins("JOIN aff_order ON aff_order.accesstrade_order_id = aff_reward.accesstrade_order_id")
		sql.Where("aff_order.campaign_id = ?", campaignId)
	}

	err := sql.Row().Scan(&cashback.Distributed, &cashback.Remain)

	if err != nil {
		return cashback, err
	}
	cashback.Remain = cashback.Remain - cashback.Distributed
	return cashback, nil
}

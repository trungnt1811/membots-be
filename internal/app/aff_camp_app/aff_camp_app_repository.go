package aff_camp_app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affCampAppRepository struct {
	db *gorm.DB
}

func NewAffCampAppRepository(db *gorm.DB) interfaces.AffCampAppRepository {
	return &affCampAppRepository{
		db: db,
	}
}

func (r affCampAppRepository) GetAllAffCampaign(ctx context.Context, orderBy string, page, size int) ([]model.AffCampaignLessApp, error) {
	var listAffCampaign []model.AffCampaignLessApp
	var err error
	offset := (page - 1) * size
	switch orderBy {
	case interfaces.ListAffCampaignOrderByMostCommission:
		orderQuery := "CASE " +
			"WHEN attribute_type = 'percent' THEN 3 " +
			"WHEN attribute_type = 'vnd' THEN 2 " +
			"ELSE 1 " +
			"END DESC, CAST(attribute_value + 0 AS DECIMAL(12,2)) DESC"
		err = r.db.Joins("Brand").
			Preload("Attributes", func(db *gorm.DB) *gorm.DB {
				db = db.Order(orderQuery)
				return db
			}).
			Where("stella_status = ?", model.StellaStatusInProgress).
			Limit(size + 1).Offset(offset).
			Find(&listAffCampaign).Error
	default:
		err = r.db.Joins("Brand").
			Preload("Attributes").
			Where("stella_status = ?", model.StellaStatusInProgress).
			Limit(size + 1).Offset(offset).
			Order("aff_campaign.id ASC").
			Find(&listAffCampaign).Error
	}
	return listAffCampaign, err
}

func (r affCampAppRepository) GetAffCampaignById(ctx context.Context, id uint64) (model.AffCampaignApp, error) {
	var affCampaign model.AffCampaignApp
	err := r.db.
		Joins("Brand").
		Preload("Attributes").
		Where("aff_campaign.id = ?", id).
		First(&affCampaign).Error
	return affCampaign, err
}

func (r affCampAppRepository) GetListAffCampaignByBrandIds(ctx context.Context, brandIds []uint, page, size int) ([]model.AffCampaignComBrand, error) {
	var listAffCampaign []model.AffCampaignComBrand
	// Ordering by the order of values in a IN() clause
	s, _ := json.Marshal(brandIds)
	findInSet := strings.Trim(string(s), "[]")
	offset := (page - 1) * size
	err := r.db.Joins("Brand").
		Preload("Attributes").
		Where("aff_campaign.brand_id IN ? AND stella_status = ?", brandIds, model.StellaStatusInProgress).
		Limit(size + 1).Offset(offset).
		Order(fmt.Sprintf("FIND_IN_SET(aff_campaign.brand_id,'%s')", findInSet)).
		Find(&listAffCampaign).Error
	return listAffCampaign, err
}

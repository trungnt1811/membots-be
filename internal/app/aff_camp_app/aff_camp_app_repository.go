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
	orderQuery := ""
	switch orderBy {
	case interfaces.ListAffCampaignOrderByMostCommission:
		orderQuery = "CASE " +
			"WHEN stella_max_com LIKE 'Upto%' THEN 2 " +
			"WHEN stella_max_com LIKE '%VND' THEN 3 " +
			"WHEN stella_max_com = '' THEN 4 " +
			"WHEN stella_max_com IS NULL THEN 5 " +
			"ELSE 1 " +
			"END ASC, CAST(REGEXP_SUBSTR(stella_max_com, '[+-]?([0-9]*[.])?[0-9]+') + 0 AS DECIMAL(12,2)) DESC"
	default:
		orderQuery = "aff_campaign.id ASC"
	}
	offset := (page - 1) * size
	err := r.db.Joins("Brand").Where("stella_status = ?", model.StellaStatusInProgress).
		Order(orderQuery).
		Find(&listAffCampaign).Limit(size + 1).Offset(offset).Error
	return listAffCampaign, err
}

func (r affCampAppRepository) GetAffCampaignById(ctx context.Context, id uint64) (model.AffCampaignApp, error) {
	var affCampaign model.AffCampaignApp
	err := r.db.Joins("Brand").Where("aff_campaign.id = ? AND stella_status = ?", id, model.StellaStatusInProgress).
		First(&affCampaign).Error
	return affCampaign, err
}

func (r affCampAppRepository) GetListAffCampaignByBrandIds(ctx context.Context, brandIds []uint64) ([]model.AffCampaignComBrand, error) {
	var listAffCampaign []model.AffCampaignComBrand
	s, _ := json.Marshal(brandIds)
	findInSet := strings.Trim(string(s), "[]")
	err := r.db.Joins("Brand").Where("aff_campaign.brand_id IN ? AND stella_status = ?", brandIds, model.StellaStatusInProgress).
		Order(fmt.Sprintf("FIND_IN_SET(aff_campaign.brand_id,'%s')", findInSet)).
		Find(&listAffCampaign).Error
	return listAffCampaign, err
}

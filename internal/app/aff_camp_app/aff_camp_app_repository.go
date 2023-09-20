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

func (r affCampAppRepository) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaignLessApp, error) {
	var listAffCampaign []model.AffCampaignLessApp
	offset := (page - 1) * size
	err := r.db.Joins("Brand").Where("stella_status = ?", model.StellaStatusInProgress).
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
	// Ordering by the order of values in a IN() clause
	s, _ := json.Marshal(brandIds)
	findInSet := strings.Trim(string(s), "[]")
	err := r.db.Joins("Brand").Where("aff_campaign.brand_id IN ? AND stella_status = ?", brandIds, model.StellaStatusInProgress).
		Order(fmt.Sprintf("FIND_IN_SET(aff_campaign.brand_id,'%s')", findInSet)).
		Find(&listAffCampaign).Error
	return listAffCampaign, err
}

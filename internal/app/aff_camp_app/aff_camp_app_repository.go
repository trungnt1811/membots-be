package aff_camp_app

import (
	"context"

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

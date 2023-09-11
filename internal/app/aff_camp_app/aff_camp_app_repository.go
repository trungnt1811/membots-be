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

func (r affCampAppRepository) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaignApp, error) {
	var listAffCampaign []model.AffCampaignApp
	offset := (page - 1) * size
	err := r.db.Find(&listAffCampaign).Limit(size + 1).Offset(offset).Error
	return listAffCampaign, err
}

func (r affCampAppRepository) GetAffCampaignById(ctx context.Context, id uint64) (model.AffCampaignApp, error) {
	var affCampaign model.AffCampaignApp
	err := r.db.Where("id = ?", id).First(&affCampaign).Error
	return affCampaign, err
}

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

func (r affCampAppRepository) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error) {
	var listAffCampaign []model.AffCampaign
	offset := (page - 1) * size
	query := "SELECT id, brand_id, accesstrade_id, created_at, updated_at, thumbnail, name, url, category_id, start_time, end_time, stella_description, stella_status, stella_max_com " +
		"FROM aff_campaign " +
		"ORDER BY id ASC LIMIT ? OFFSET ?"
	err := r.db.Raw(query, size+1, offset).Scan(&listAffCampaign).Error
	return listAffCampaign, err
}

func (r affCampAppRepository) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error) {
	var affCampaign model.AffCampaign
	query := "SELECT id, brand_id, accesstrade_id, created_at, updated_at, thumbnail, name, url, category_id, start_time, end_time, stella_description, stella_status, stella_max_com " +
		"FROM aff_campaign " +
		"WHERE accesstrade_id = ?"
	err := r.db.Raw(query, accesstradeId).Scan(&affCampaign).Error
	return affCampaign, err
}

package mysql

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affCampaignRepository struct {
	db *gorm.DB
}

func (r affCampaignRepository) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error) {
	var listAffCampaign []model.AffCampaign
	offset := (page - 1) * size
	query := "SELECT stella_description, name, thumbnail, url, max_com, stella_max_com, start_time, end_time, category_id, brand_id, stella_status " +
		"FROM aff_campaign " +
		"ORDER BY id ASC LIMIT ? OFFSET ?"
	err := r.db.Raw(query, size+1, offset).Scan(&listAffCampaign).Error
	return listAffCampaign, err
}

func (r affCampaignRepository) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error) {
	var affCampaign model.AffCampaign
	query := "SELECT stella_description, name, thumbnail, url, max_com, stella_max_com, start_time, end_time, category_id, brand_id, stella_status " +
		"FROM aff_campaign " +
		"WHERE accesstrade_id = ?"
	err := r.db.Raw(query, accesstradeId).Scan(&affCampaign).Error
	return affCampaign, err
}

func NewAffCampaignRepository(db *gorm.DB) interfaces.AffCampaignRepository {
	return &affCampaignRepository{
		db: db,
	}
}

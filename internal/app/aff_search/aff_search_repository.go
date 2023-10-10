package aff_search

import (
	"context"
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affBannerRepository struct {
	Db *gorm.DB
}

func (a *affBannerRepository) Search(ctx context.Context, q string, page, size int) (model.AffSearch, error) {
	var searchResult model.AffSearch
	var totalAffCampaign int64
	var listAffCampaign []model.AffCampaignLessApp
	offset := (page - 1) * size

	err := a.Db.Table("aff_campaign").
		Joins("Brand").
		Preload("Attributes").
		Where("aff_campaign.stella_status = ? AND (MATCH (aff_campaign.name) AGAINST(? IN NATURAL LANGUAGE MODE) OR "+
			"MATCH (Brand.name) AGAINST(? IN NATURAL LANGUAGE MODE))", "IN_PROGRESS", fmt.Sprint(q, " *"), fmt.Sprint(q, " *")).
		Limit(size + 1).
		Offset(offset).
		Find(&listAffCampaign).Error
	if err != nil {
		return searchResult, err
	}
	err = a.Db.Raw("select count(*) FROM aff_campaign as ac LEFT JOIN brand as b ON b.id = ac.brand_id "+
		"WHERE ac.stella_status = ? AND (MATCH (ac.name) AGAINST(? IN NATURAL LANGUAGE MODE) OR "+
		"MATCH (b.name) AGAINST(? IN NATURAL LANGUAGE MODE))",
		"IN_PROGRESS", fmt.Sprint(q, " *"), fmt.Sprint(q, " *")).Scan(&totalAffCampaign).Error
	if err != nil {
		return searchResult, err
	}
	return model.AffSearch{
		AffCampaign:   listAffCampaign,
		TotalCampaign: totalAffCampaign,
	}, nil
}

func NewAffSearchRepository(db *gorm.DB) interfaces.AffSearchRepository {
	return &affBannerRepository{
		Db: db,
	}
}

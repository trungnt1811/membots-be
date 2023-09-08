package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affCampaignRepository struct {
	Db *gorm.DB
}

func (a *affCampaignRepository) GetCampaignById(id uint) (model.AffCampaign, error) {
	var affCampaign model.AffCampaign
	if err := a.Db.Table("aff_campaign").
		Joins("Description").
		Where("aff_campaign.id = ?", id).First(&affCampaign).Error; err != nil {
		return affCampaign, err
	}
	return affCampaign, nil
}

func (a *affCampaignRepository) UpdateCampaign(id uint, updates map[string]interface{}) error {
	return a.Db.Table("aff_campaign").Where("id = ?", id).Updates(updates).Error
}

func (a *affCampaignRepository) GetAllCampaign(listStatus []string, page, size int) ([]model.AffCampaign, error) {
	var listAffCampaign []model.AffCampaign
	offset := (page - 1) * size
	if err := a.Db.Table("aff_campaign").
		Joins("Description").
		Where("aff_campaign.stella_status IN ?", listStatus).
		Limit(size + 1).
		Offset(offset).
		Find(&listAffCampaign).
		Error; err != nil {
		return listAffCampaign, err
	}
	return listAffCampaign, nil
}

func (a *affCampaignRepository) CountCampaign(listStatus []string) (int64, error) {
	var total int64
	if err := a.Db.Table("aff_campaign").Where("stella_status IN ?", listStatus).Count(&total).Error; err != nil {
		return total, err
	}
	return total, nil
}

func NewConsoleCampaignRepository(db *gorm.DB) interfaces.ConsoleCampRepository {
	return &affCampaignRepository{
		Db: db,
	}
}

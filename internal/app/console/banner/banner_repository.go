package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affBannerRepository struct {
	Db *gorm.DB
}

func (a *affBannerRepository) GetBannerById(id uint) (model.AffBanner, error) {
	var affBanner model.AffBanner
	if err := a.Db.Table("aff_banner").
		Joins("AffCampaign").
		Where("aff_banner.id = ?", id).First(&affBanner).Error; err != nil {
		return affBanner, err
	}
	return affBanner, nil
}

func (a *affBannerRepository) UpdateBanner(id uint, updates map[string]interface{}) error {
	return a.Db.Table("aff_banner").Where("id = ?", id).Updates(updates).Error
}

func (a *affBannerRepository) GetAllBanner(listStatus []string, page, size int) ([]model.AffBanner, error) {
	var listAffBanner []model.AffBanner
	offset := (page - 1) * size
	if err := a.Db.Table("aff_banner").
		Joins("AffCampaign").
		Where("aff_banner.status IN ?", listStatus).
		Limit(size + 1).
		Offset(offset).
		Find(&listAffBanner).
		Error; err != nil {
		return listAffBanner, err
	}
	return listAffBanner, nil
}

func (a *affBannerRepository) CountBanner(listStatus []string) (int64, error) {
	var total int64
	if err := a.Db.Table("aff_banner").Where("status IN ?", listStatus).Count(&total).Error; err != nil {
		return total, err
	}
	return total, nil
}

func NewConsoleBannerRepository(db *gorm.DB) interfaces.ConsoleBannerRepository {
	return &affBannerRepository{
		Db: db,
	}
}

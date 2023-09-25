package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type ConsoleCampRepository interface {
	GetAllCampaign(status []string, q string, page, size int) ([]model.AffCampaignLessApp, error)
	GetCampaignById(id uint) (model.AffCampaign, error)
	GetCampaignLessByAccessTradeId(accessTradeId string) (model.AffCampaignLess, error)
	CountCampaign(status []string, q string) (int64, error)
	UpdateCampaign(id uint, updates map[string]interface{}) error
	UpdateCampaignAttribute(id uint, attributes []model.AffCampaignAttribute) error
	LoadCategoryFromBrandId(brandId uint) (uint, error)
}

type ConsoleCampUCase interface {
	GetAllCampaign(status []string, q string, page, size int) (dto.AffCampaignAppDtoResponse, error)
	UpdateCampaign(id uint, campaign dto.AffCampaignAppDto) error
	GetCampaignById(id uint) (dto.AffCampaignDto, error)
}

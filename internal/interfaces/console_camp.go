package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type ConsoleCampRepository interface {
	GetAllCampaign(status []string, page, size int) ([]model.AffCampaignLessApp, error)
	GetCampaignById(id uint) (model.AffCampaign, error)
	GetCampaignLessByAccessTradeId(accessTradeId string) (model.AffCampaignLess, error)
	CountCampaign(status []string) (int64, error)
	UpdateCampaign(id uint, updates map[string]interface{}) error
}

type ConsoleCampUCase interface {
	GetAllCampaign(status []string, page, size int) (dto.AffCampaignAppDtoResponse, error)
	UpdateCampaign(id uint, campaign dto.AffCampaignAppDto) error
	GetCampaignById(id uint) (dto.AffCampaignDto, error)
}

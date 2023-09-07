package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type ConsoleCampRepository interface {
	GetAllCampaign(status []string, page, size int) ([]model.AffCampaign, error)
	CountCampaign(status []string) (int64, error)
	UpdateCampaign(id uint, updates map[string]interface{}) error
}

type ConsoleCampUCase interface {
	GetAllCampaign(status []string, page, size int) (dto.AffCampaignDtoResponse, error)
	UpdateCampaign(id uint, campaign dto.AffCampaignDto) error
}

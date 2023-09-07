package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/dto"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
)

type ConsoleCampRepository interface {
	GetAllCampaign(status []string, page, size int) ([]model.AffCampaign, error)
	CountCampaign(status []string) (int64, error)
}

type ConsoleCampUCase interface {
	GetAllCampaign(status []string, page, size int) (dto.AffCampaignDtoResponse, error)
}

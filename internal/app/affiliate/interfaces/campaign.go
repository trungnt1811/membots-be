package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/dto"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
)

type CampaignRepository interface {
	RetrieveCampaigns(q map[string]any) ([]model.AffCampaign, error)
	RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]*model.AffCampaign, error)
	SaveATCampaign(atCampaign *types.ATCampaign) error

	CreateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error)
	UpdateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error)
	DeactivateCampaigns(data []model.AffCampaign) error

	RetrieveAffLinks(campaignId uint) ([]model.AffLink, error)
	CreateAffLinks(data []model.AffLink) error
}

type CampaignUsecase interface {
	GenerateAffLink(*model.UserEntity, *dto.CreateLinkPayload) (*dto.CreateLinkResponse, error)
}

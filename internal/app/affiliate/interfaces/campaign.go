package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/dto"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
)

type CampaignRepository interface {
	RetrieveCampaigns(q map[string]any) ([]model.Campaign, error)
	RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]*model.Campaign, error)
	SaveATCampaign(atCampaign *types.ATCampaign, atMerchant *types.ATMerchant) error

	CreateCampaigns(data []model.Campaign) ([]model.Campaign, error)
	UpdateCampaigns(data []model.Campaign) ([]model.Campaign, error)
	DeactivateCampaigns(data []model.Campaign) error

	RetrieveAffLinks(campaignId uint) ([]model.AffLink, error)
	CreateAffLinks(data []model.AffLink) error
}

type CampaignUsecase interface {
	GenerateAffLink(*model.UserEntity, *dto.CreateLinkPayload) (*dto.CreateLinkResponse, error)
}

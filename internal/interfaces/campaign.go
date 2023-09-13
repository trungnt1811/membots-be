package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	model2 "github.com/astraprotocol/affiliate-system/internal/model"
)

type CampaignRepository interface {
	RetrieveCampaigns(q map[string]any) ([]model2.AffCampaign, error)
	RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]*model2.AffCampaign, error)
	SaveATCampaign(atCampaign *types.ATCampaign) error

	CreateCampaigns(data []model2.AffCampaign) ([]model2.AffCampaign, error)
	UpdateCampaigns(data []model2.AffCampaign) ([]model2.AffCampaign, error)
	DeactivateCampaigns(data []model2.AffCampaign) error

	RetrieveAffLinks(campaignId uint) ([]model2.AffLink, error)
	CreateAffLinks(data []model2.AffLink) error
}

type CampaignUCase interface {
	GenerateAffLink(uint64, *dto.CreateLinkPayload) (*dto.CreateLinkResponse, error)
}

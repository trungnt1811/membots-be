package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
	model2 "github.com/astraprotocol/affiliate-system/internal/model"
)

type CampaignRepository interface {
	GetCampaignLessById(campaignId uint) (model2.AffCampaignLess, error)
	RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]model2.AffCampaign, error)
	SaveATCampaign(atCampaign *types.ATCampaign) (*model2.AffCampaign, error)
	QueryActiveIds() ([]uint, error)
	DeactivateCampaign(campaignId uint) (*model2.AffCampaign, error)
	DeactivateCampaignLinks(campaignId uint) error

	CreateCampaigns(data []model2.AffCampaign) ([]model2.AffCampaign, error)
	UpdateCampaigns(data []model2.AffCampaign) ([]model2.AffCampaign, error)
	UpdateCampaignByID(ID uint, updates map[string]any, description map[string]any) error

	RetrieveAffLinks(campaignId uint, originalUrl string) ([]model2.AffLink, error)
	CreateAffLinks(data []model2.AffLink) error

	CreateTrackedClick(*model2.AffTrackedClick) error
}

type CampaignUCase interface {
	GenerateAffLink(uint64, *dto.CreateLinkPayload) (*dto.CreateLinkResponse, error)
}

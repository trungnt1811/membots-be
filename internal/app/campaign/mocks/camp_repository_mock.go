package mocks

import (
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type CampaignRepoMock struct {
	campaigns []model.AffCampaign
}

func NewCampaignRepoMock() *CampaignRepoMock {
	return &CampaignRepoMock{}
}

func (repo *CampaignRepoMock) GetCampaignById(q map[string]any) ([]model.AffCampaign, error) {
	return []model.AffCampaign{
		{
			ID:            14,
			AccessTradeId: "4751584435713464237",
			Name:          "Shopee Việt Nam Smartlink cho tất cả thiết bị ",
		},
	}, nil
}

func (repo *CampaignRepoMock) GetCampaignLessById(campaignId uint) (model.AffCampaignLess, error) {
	return model.AffCampaignLess{}, nil
}

func (repo *CampaignRepoMock) RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]*model.AffCampaign, error) {
	return nil, nil
}
func (repo *CampaignRepoMock) SaveATCampaign(atCampaign *types.ATCampaign) error {
	return nil
}

func (repo *CampaignRepoMock) CreateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error) {
	return nil, nil
}
func (repo *CampaignRepoMock) UpdateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error) {
	return nil, nil
}
func (repo *CampaignRepoMock) DeactivateCampaigns(data []model.AffCampaign) error {
	return nil
}

func (repo *CampaignRepoMock) RetrieveAffLinks(campaignId uint, originalUrl string) ([]model.AffLink, error) {
	return []model.AffLink{}, nil
}
func (repo *CampaignRepoMock) CreateAffLinks(data []model.AffLink) error {
	return nil
}
func (repo *CampaignRepoMock) CreateTrackedClick(*model.AffTrackedClick) error {
	return nil
}

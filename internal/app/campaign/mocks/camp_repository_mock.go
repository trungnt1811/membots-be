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

func (repo *CampaignRepoMock) QueryActiveIds() ([]uint, error) {
	return []uint{}, nil
}

func (repo *CampaignRepoMock) DeactivateCampaign(campaignId uint) (*model.AffCampaign, error) {
	return &model.AffCampaign{
		ID:            14,
		AccessTradeId: "4751584435713464237",
		Name:          "Shopee Việt Nam Smartlink cho tất cả thiết bị ",
		StellaStatus:  model.StellaStatusEnded,
	}, nil
}

func (repo *CampaignRepoMock) GetCampaignLessById(campaignId uint) (model.AffCampaignLess, error) {
	return model.AffCampaignLess{}, nil
}

func (repo *CampaignRepoMock) RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]model.AffCampaign, error) {
	return nil, nil
}
func (repo *CampaignRepoMock) SaveATCampaign(atCampaign *types.ATCampaign) (*model.AffCampaign, error) {
	return nil, nil
}

func (repo *CampaignRepoMock) CreateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error) {
	return nil, nil
}
func (repo *CampaignRepoMock) UpdateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error) {
	return nil, nil
}
func (repo *CampaignRepoMock) UpdateCampaignByID(ID uint, updates map[string]any, description map[string]any) error {
	return nil
}

func (repo *CampaignRepoMock) DeactivateCampaignLinks(campaignId uint) error {
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

package interfaces

import "github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade/types"

type ATRepository interface {
	QueryMerchants() ([]types.ATMerchant, error)
	QueryCampaigns(onlyApproval bool, page int, limit int) (*types.ATCampaignListResp, error)
	QueryOrders(q types.ATOrderQuery, page int, limit int) (*types.ATOrderListResp, error)
	CreateTrackingLinks(campaignId string, urls []string, utm map[string]string) (*types.ATLinkResp, error)
}

type ATUsecase interface {
	QueryAndSaveCampaigns(onlyApproval bool) (int, error)
	CreateAndSaveLink() (int, error)
}

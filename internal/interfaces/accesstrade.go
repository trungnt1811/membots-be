package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
)

type ATRepository interface {
	QueryMerchants() ([]types.ATMerchant, error)
	QueryCampaigns(onlyApproval bool, page int, limit int) (*types.ATCampaignListResp, error)
	QueryTransactions(q types.ATTransactionQuery, page int, limit int) (*types.ATTransactionResp, error)
	QueryOrders(q types.ATOrderQuery, page int, limit int) (*types.ATOrderListResp, error)
	CreateTrackingLinks(campaignId string, urls []string, utm map[string]string) (*types.ATLinkResp, error)
}

type ATUCase interface {
	QueryAndSaveCampaigns(onlyApproval bool) (int, error)
	CreateAndSaveLink() (int, error)
}

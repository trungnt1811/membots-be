package interfaces

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
)

type ATRepository interface {
	QueryMerchants() ([]types.ATMerchant, error)
	QueryCampaigns(onlyApproval bool, page int, limit int) (time.Time, *types.ATCampaignListResp, error)
	QueryTransactions(q types.ATTransactionQuery, page int, limit int) (*types.ATTransactionResp, error)
	QueryOrders(q types.ATOrderQuery, page int, limit int) (*types.ATOrderListResp, error)
	CreateTrackingLinks(campaignId string, shorten bool, urls []string, utm map[string]string) (*types.ATLinkResp, error)
}

type ATUCase interface {
	QueryAndSaveCampaigns(onlyApproval bool) (int, int, error)
}

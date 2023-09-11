package mocks

import (
	"encoding/json"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
)

type AccessTradeRepoMock struct{}

func NewAccessTradeRepoMock() *AccessTradeRepoMock {
	return &AccessTradeRepoMock{}
}

func (repo *AccessTradeRepoMock) QueryMerchants() ([]types.ATMerchant, error) {
	return nil, nil
}
func (repo *AccessTradeRepoMock) QueryCampaigns(onlyApproval bool, page int, limit int) (*types.ATCampaignListResp, error) {
	return nil, nil
}
func (repo *AccessTradeRepoMock) QueryTransactions(q types.ATTransactionQuery, page int, limit int) (*types.ATTransactionResp, error) {
	var resp types.ATTransactionResp
	err := json.Unmarshal([]byte(SAMPLE_TXS_RESP), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (repo *AccessTradeRepoMock) QueryOrders(q types.ATOrderQuery, page int, limit int) (*types.ATOrderListResp, error) {
	var resp types.ATOrderListResp
	err := json.Unmarshal([]byte(SAMPLE_ORDER_LIST_RESP), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (repo *AccessTradeRepoMock) CreateTrackingLinks(campaignId string, urls []string, utm map[string]string) (*types.ATLinkResp, error) {
	return nil, nil
}

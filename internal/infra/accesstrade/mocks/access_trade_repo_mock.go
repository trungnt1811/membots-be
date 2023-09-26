package mocks

import (
	"encoding/json"

	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
)

type AccessTradeRepoMock struct{}

func NewAccessTradeRepoMock() *AccessTradeRepoMock {
	return &AccessTradeRepoMock{}
}

func (repo *AccessTradeRepoMock) QueryMerchants() ([]types.ATMerchant, error) {
	return nil, nil
}
func (repo *AccessTradeRepoMock) QueryCampaigns(onlyApproval bool, page int, limit int) (*types.ATCampaignListResp, error) {
	var resp types.ATCampaignListResp
	err := json.Unmarshal([]byte(SAMPLE_CAMPAIGNS), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
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
func (repo *AccessTradeRepoMock) CreateTrackingLinks(campaignId string, shorten bool, urls []string, utm map[string]string) (*types.ATLinkResp, error) {
	return &types.ATLinkResp{
		Data: types.AllLink{
			ErrorLink: []types.ErrorLink{},
			SuccessLink: []types.ATLink{
				{
					AffLink:   "https://go.isclix.com/deep_link/6243337751671016280/4751584435713464237?url=https%3A%2F%2Fshopee.vn%2F\u0026sub5=pub-api\u0026utm_source=stella",
					FirstLink: "",
					ShortLink: "https://shorten.asia/wAVqMD4U",
					UrlOrigin: "",
				},
			},
			SuspendUrl: []any{},
		},
		Success: true,
	}, nil
}

package accesstrade

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"

	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/util/log"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	AccesstradeEndpoint     = "https://api.accesstrade.vn/v1"
	SERVICE_QUERY_CAMPAIGNS = "query-campaigns"
	SERVICE_QUERY_TXS       = "query-txs"
	SERVICE_QUERY_MERCHANTS = "query-merchants"
	SERVICE_QUERY_ORDERS    = "query-orders"
)

var serviceKeys = []string{SERVICE_QUERY_CAMPAIGNS, SERVICE_QUERY_TXS, SERVICE_QUERY_MERCHANTS, SERVICE_QUERY_ORDERS}

type accessTradeRepository struct {
	APIKey       string
	tickLimiters map[string]<-chan time.Time
	caller       *resty.Client
}

func NewAccessTradeRepository(APIKey string, retry int, timeoutSec int, tickLimiterMs int) interfaces.ATRepository {
	if APIKey == "" {
		panic(errors.New("no AccessTrade API Key"))
	}
	// Initialize and configure resty client
	client := resty.New()
	client.SetRetryCount(retry)
	client.SetTimeout(time.Duration(timeoutSec * int(time.Second)))

	// For waiting all the requests to be done before start a new one
	ticks := map[string]<-chan time.Time{}
	for _, k := range serviceKeys {
		ticks[k] = time.Tick(time.Duration(tickLimiterMs) * time.Millisecond)
	}

	return &accessTradeRepository{
		APIKey:       APIKey,
		tickLimiters: ticks,
		caller:       client,
	}
}

func (r *accessTradeRepository) initWithHeaders() *resty.Request {
	req := r.caller.R()
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", fmt.Sprintf("Token %s", r.APIKey))
	return req
}

func (r *accessTradeRepository) QueryMerchants() ([]types.ATMerchant, error) {
	url := fmt.Sprintf("%s/offers_informations/merchant_list", AccesstradeEndpoint)
	req := r.initWithHeaders()
	fmt.Println("QueryMerchants", url)

	<-r.tickLimiters[SERVICE_QUERY_MERCHANTS]
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	var body types.ATMerchantResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, err
	}
	return body.Data, nil
}

func (r *accessTradeRepository) QueryCampaigns(onlyApproval bool, page int, limit int) (time.Time, *types.ATCampaignListResp, error) {
	url := fmt.Sprintf("%s/campaigns", AccesstradeEndpoint)
	req := r.initWithHeaders()

	if onlyApproval {
		req.SetQueryParam("approval", "successful")
	}
	if page != 0 {
		req.SetQueryParam("page", fmt.Sprint(page))
	}
	if limit != 0 {
		req.SetQueryParam("limit", fmt.Sprint(limit))
	}

	// Wait for tickLimiter to return before start new request
	t := <-r.tickLimiters[SERVICE_QUERY_CAMPAIGNS]
	resp, err := req.Get(url)
	if err != nil {
		return t, nil, err
	}

	var body types.ATCampaignListResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return t, nil, err
	}
	return t, &body, nil
}

// The `QueryTransactions` function is used to query transactions from the AccessTrade API. It takes in parameters
// such as `q` (of type `types.ATTransactionQuery`), `page`, and `limit`.
func (r *accessTradeRepository) QueryTransactions(q types.ATTransactionQuery, page int, limit int) (*types.ATTransactionResp, error) {
	url := fmt.Sprintf("%s/transactions", AccesstradeEndpoint)
	req := r.initWithHeaders()

	if page != 0 {
		req.SetQueryParam("page", fmt.Sprint(page))
	}
	if limit != 0 {
		req.SetQueryParam("limit", fmt.Sprint(limit))
	}

	params, err := query.Values(&q)
	if err != nil {
		return nil, err
	}
	req.SetQueryString(params.Encode())

	// Start the connection
	// Wait for tickLimiter to return before start new request
	<-r.tickLimiters[SERVICE_QUERY_TXS]
	resp, err := req.Get(url)
	if err != nil {
		return nil, errors.Errorf("request error: %v", err)
	}

	// Parse response body
	var body types.ATTransactionResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, errors.Errorf("parse json resp error: %v", err)
	}
	return &body, nil
}

// QueryOrders The `QueryOrders` function is used to query orders from the AccessTrade API. It takes in parameters
// such as `q` (of type `types.ATOrderQuery`), `page`, and `limit`.
func (r *accessTradeRepository) QueryOrders(q types.ATOrderQuery, page int, limit int) (*types.ATOrderListResp, error) {
	url := fmt.Sprintf("%s/order-list", AccesstradeEndpoint)
	req := r.initWithHeaders()

	if page != 0 {
		req.SetQueryParam("page", fmt.Sprint(page))
	}
	if limit != 0 {
		req.SetQueryParam("limit", fmt.Sprint(limit))
	}

	params, err := query.Values(&q)
	if err != nil {
		return nil, err
	}
	req.SetQueryString(params.Encode())

	// Start the connection
	// Wait for tickLimiter to return before start new request
	<-r.tickLimiters[SERVICE_QUERY_ORDERS]
	resp, err := req.Get(url)
	if err != nil {
		return nil, errors.Errorf("request error: %v", err)
	}

	if !resp.IsSuccess() {
		errMsg := string(resp.Body())
		return nil, errors.Errorf("response error: %d - %s", resp.StatusCode(), errMsg)
	}

	// Parse response body
	var body types.ATOrderListResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, errors.Errorf("parse json resp error: %v", err)
	}
	return &body, nil
}

// CreateTrackingLinks The `CreateTrackingLinks` function is used to create tracking links for a campaign in the
// AccessTrade API. It takes in parameters such as `campaignId` (string), `urls` (slice of strings),
// and `additional` (map[string]string).
func (r *accessTradeRepository) CreateTrackingLinks(campaignId string, shorten bool, urls []string, additional map[string]string) (*types.ATLinkResp, error) {
	url := fmt.Sprintf("%s/product_link/create", AccesstradeEndpoint)
	req := r.initWithHeaders()

	reqBody := map[string]any{
		"campaign_id":    campaignId,
		"urls":           urls,
		"create_shorten": shorten,
		"url_enc":        true,
		"utm_source":     "stella",
	}
	for k, v := range additional {
		reqBody[k] = v
	}
	req.SetBody(reqBody)

	// Start the connection
	resp, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	// Parse response body
	var body types.ATLinkResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		// omit some marshal failed
		log.LG.Errorf("marshal response failed: %v", err)
	}
	if !body.Success {
		if len(body.Data.ErrorLink) != 0 {
			// Loop through error link and response
			for _, link := range body.Data.ErrorLink {
				if link.Message != "" {
					return nil, fmt.Errorf("%s: %s", link.Message, link.UrlOrigin)
				}
			}
		}
		return nil, fmt.Errorf("failed status: %d", body.StatusCode)
	}
	return &body, nil
}

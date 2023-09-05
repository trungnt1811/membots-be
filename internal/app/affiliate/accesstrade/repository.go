package accesstrade

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade/types"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	ACCESSTRADE_ENDPOINT = "https://api.accesstrade.vn/v1"
)

type AccessTradeRepository struct {
	APIKey string
	caller *resty.Client
}

func NewAccessTradeRepository(APIKey string, retry int, timeoutSec int) *AccessTradeRepository {
	if APIKey == "" {
		// TODO: Alert about required API key
	}
	// Initialize and configure resty client
	client := resty.New()
	client.SetRetryCount(retry)
	client.SetTimeout(time.Duration(timeoutSec * int(time.Second)))

	return &AccessTradeRepository{
		APIKey: APIKey,
		caller: client,
	}
}

func (r *AccessTradeRepository) initWithHeaders() *resty.Request {
	req := r.caller.R()
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", fmt.Sprintf("Token %s", r.APIKey))
	return req
}

func (r *AccessTradeRepository) QueryMerchants() ([]types.ATMerchant, error) {
	url := fmt.Sprintf("%s/offers_informations/merchant_list", ACCESSTRADE_ENDPOINT)
	req := r.initWithHeaders()

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

func (r *AccessTradeRepository) QueryCampaigns(onlyApproval bool, page int, limit int) (*types.ATCampaignListResp, error) {
	url := fmt.Sprintf("%s/campaigns", ACCESSTRADE_ENDPOINT)
	req := r.initWithHeaders()

	if onlyApproval {
		req.SetQueryParam("approval", "successful")
	}
	req.SetQueryParams(map[string]string{
		"page":  fmt.Sprint(page),
		"limit": fmt.Sprint(limit),
	})

	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	var body types.ATCampaignListResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// The `QueryOrders` function is used to query orders from the AccessTrade API. It takes in parameters
// such as `q` (of type `types.OrderQuery`), `page`, and `limit`.
func (r *AccessTradeRepository) QueryOrders(q types.ATOrderQuery, page int, limit int) (*types.ATOrderListResp, error) {
	url := fmt.Sprintf("%s/order-list", ACCESSTRADE_ENDPOINT)
	req := r.initWithHeaders()

	req.SetQueryParams(map[string]string{
		"page":  fmt.Sprint(page),
		"limit": fmt.Sprint(limit),
	})

	params, err := query.Values(&q)
	if err != nil {
		return nil, err
	}
	req.SetQueryString(params.Encode())

	// Start the connection
	resp, err := req.Get(url)
	if err != nil {
		return nil, errors.Errorf("request error: %v", err)
	}

	// Parse response body
	var body types.ATOrderListResp
	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		return nil, errors.Errorf("parse json resp error: %v", err)
	}
	return &body, nil
}

// The `CreateTrackingLinks` function is used to create tracking links for a campaign in the
// AccessTrade API. It takes in parameters such as `campaignId` (string), `urls` (slice of strings),
// and `additional` (map[string]string).
func (r *AccessTradeRepository) CreateTrackingLinks(campaignId string, urls []string, additional map[string]string) (*types.ATLinkResp, error) {
	url := fmt.Sprintf("%s/product_link/create", ACCESSTRADE_ENDPOINT)
	req := r.initWithHeaders()

	reqBody := map[string]any{
		"campaign_id":    campaignId,
		"urls":           urls,
		"create_shorten": true,
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
		return nil, err
	}
	return &body, nil
}

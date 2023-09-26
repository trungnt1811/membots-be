package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

const (
	THRESHOLD_VOLUME      = 50000 // Ensure if user who have under 50k asa cannot manipulate the price
	CURRENT_PRICE         = 0
	FOR_SEND_TO_USER      = 1
	FOR_RECEIVE_FROM_USER = 2
)

type TikiClientConfig struct {
	BaseUrl string
	ApiKey  string
}

type TikiClient struct {
	*req.Client
	config TikiClientConfig
}

func NewTikiClient(config TikiClientConfig) *TikiClient {
	client := req.C().
		SetUserAgent("stella-affiliate-backend"). // Chainable client settings.
		SetBaseURL(config.BaseUrl).
		SetCommonBearerAuthToken(config.ApiKey).
		SetTimeout(3 * time.Second).
		SetCommonRetryCount(3).
		SetCommonErrorResult(&dto.ErrorMessage{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
				return nil
			}
			if errMsg, ok := resp.ErrorResult().(*dto.ErrorMessage); ok {
				resp.Err = errMsg // Convert api error into go error
				return nil
			}
			if !resp.IsSuccessState() {
				// Neither a success response nor a error response, record details to help troubleshooting
				resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
			}
			return nil
		})
	return &TikiClient{
		Client: client,
		config: config,
	}
}

func (client *TikiClient) GetUrl(url string) string {
	return client.config.BaseUrl + url
}

func (client *TikiClient) GetAstraPrice(ctx context.Context) (int64, error) {
	return client.GetAstraPriceFromExchange(ctx, FOR_SEND_TO_USER)
}

// GetAstraPrice get ASA price from Tiki exchange
// @Param forSendReward - is the ASA price calculated for send to customer (true) or in opposite way (false)
// - true: calculate base on asks price, to prevent user push the price down
// - false: calculate base on bids price, to prevent user push the price up
func (c *TikiClient) GetAstraPriceFromExchange(ctx context.Context, calculateMode int) (int64, error) {
	endpoint := "/sandseel/api/v2/public/markets/asaxu/depth"

	var errRes ErrorResponseData
	var priceRes ExchangePriceResponse

	resp, err := c.R().
		SetErrorResult(&errRes).
		SetSuccessResult(&priceRes).
		Get(endpoint)
	if err != nil {
		return 0, err
	}
	if !resp.IsSuccessState() {
		log.Error().Msgf("bad response status: %v", resp.Status)
		return 0, nil
	}

	exchangePrice := priceRes.ToExchangePrice()
	price := exchangePrice.GetPrice(calculateMode, float64(THRESHOLD_VOLUME))
	return int64(price), nil
}

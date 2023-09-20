package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/imroc/req/v3"
)

type TikiClientConfig struct {
	BaseUrl        string
	ApiKey         string
	ExchangeApiKey string
}

type TikiClient struct {
	*req.Client
	config TikiClientConfig
}

func NewTikiClient(config TikiClientConfig) *TikiClient {
	client := req.C().
		SetUserAgent("reward-app-backend"). // Chainable client settings.
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
	return 200, nil
}

// func (c *TikiClient) LinkUserAccount(payload *LinkAccountPayload) (LinkAccountResponse, ErrorResponse, error) {
// 	endpoint := "/charon/partners/link/verify"

// 	var errRes ErrorResponseData
// 	var linkResponse LinkAccountResponse

// 	resp, err := c.R().SetBody(payload).
// 		SetErrorResult(&errRes).
// 		SetSuccessResult(&linkResponse).
// 		Post(endpoint)
// 	if err != nil {
// 		return LinkAccountResponse{}, errRes.Error, err
// 	}
// 	if !resp.IsSuccessState() {
// 		log.Error().Msgf("bad response status: %v", resp.Status)
// 		return LinkAccountResponse{}, errRes.Error, nil
// 	}

// 	return linkResponse, errRes.Error, nil
// }

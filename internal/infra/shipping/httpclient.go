package shipping

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

type ShippingClientConfig struct {
	BaseUrl string
	ApiKey  string
}

type ShippingClient struct {
	*req.Client
	config ShippingClientConfig
}

func NewShippingClient(config ShippingClientConfig) *ShippingClient {
	client := req.C().
		SetUserAgent("astra-affiliate-backend").
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
	return &ShippingClient{
		Client: client,
		config: config,
	}
}

func (client *ShippingClient) GetUrl(url string) string {
	return client.config.BaseUrl + url
}

func (c *ShippingClient) SendReward(payload *ReqSendPayload) (ReqSendResponse, error) {
	endpoint := "/api/v1/reward/send"

	var errRes string
	var response ReqSendResponse

	resp, err := c.R().SetBody(payload).
		SetErrorResult(&errRes).
		SetSuccessResult(&response).
		Post(endpoint)
	if err != nil {
		return ReqSendResponse{}, err
	}

	if !resp.IsSuccessState() {
		log.Error().Msgf("request send reward: bad response status: %v", resp.Status)
		return ReqSendResponse{}, nil
	}

	return response, nil
}

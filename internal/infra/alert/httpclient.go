package alert

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
)

type AlertClient struct {
	*req.Client
}

func NewAlertClient(webhookUrl string) *AlertClient {
	client := req.C().
		SetUserAgent("astra-affiliate-backend").
		SetBaseURL(webhookUrl).
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
	return &AlertClient{
		Client: client,
	}
}

func (c *AlertClient) SendMessage(message string) error {
	var errRes interface{}
	payload := &SendPayload{
		Content: message,
	}

	resp, err := c.R().SetBody(payload).
		SetErrorResult(&errRes).
		Post("")
	if err != nil {
		log.LG.Errorf("send alert failed with error %v", errRes)
		return errors.Wrapf(err, "send alert failed")
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("send alert failed: bad response status: %v", resp.Status)
	}

	return nil
}

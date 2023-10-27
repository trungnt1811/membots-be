package discord

import (
	"errors"
	"fmt"

	"github.com/astraprotocol/affiliate-system/internal/util"
)

type DiscordSender struct {
	DiscordWebhook string
}

func (sender *DiscordSender) SendMsg(topic, msg string) error {
	if sender.DiscordWebhook == "" {
		return errors.New("no discord webhook")
	}
	resp, err := util.NewHttpRequestBuilder().SetBody(map[string]any{
		"content": fmt.Sprintf("========== [AFFILIATE SYSTEM - %s] ==========\n%s", topic, msg),
	}).Build().Post(sender.DiscordWebhook)
	if err != nil {
		return err
	}
	if resp.IsSuccess() {
		return nil
	}
	return errors.New("send discord msg receive error")
}

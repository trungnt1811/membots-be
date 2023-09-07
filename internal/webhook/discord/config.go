package discord

import (
	"fmt"
)

// DiscordConfig consists of configurations of the DiscordHook.
type DiscordConfig struct {
	// Token is the OAuth token for pushing messages.
	Token string `json:"token"`

	// SubChannels is the list of channels for posting messages.
	//
	// Each channel must start with an `@` symbol. For example: @astra_alert
	SubChannels []string `json:"subChannels"`

	// MessageQueueSize specifies the message queue size.
	MessageQueueSize int `json:"messageQueueSize"`
}

func DefaultConfig() DiscordConfig {
	return DiscordConfig{
		Token:            "DISCORD_TOKEN",
		SubChannels:      []string{"SUB_CHANNEL_ID"},
		MessageQueueSize: 1024,
	}
}

// IsValid checks if the current TeleConfig is valid.
func (cfg DiscordConfig) IsValid() (bool, error) {
	if cfg.Token == "" {
		return false, fmt.Errorf("empty bot token")
	}

	if len(cfg.SubChannels) == 0 {
		return false, fmt.Errorf("empty subscribing channel")
	}

	if cfg.MessageQueueSize == 0 {
		return false, fmt.Errorf("MessageQueueSize must be greater than 0")
	}

	return true, nil
}

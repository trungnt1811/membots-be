package discord

import (
	"testing"
	"time"
)

var hook *DiscordHook

func init() {
	cfg := DefaultConfig()
	// cfg.Token =
	// cfg.SubChannels = []string{""}
	var err error
	hook, err = NewDiscordHook(cfg)
	if err != nil {
		panic(err)
	}

	err = hook.Start()
	if err != nil {
		panic(err)
	}
}

func sleepAndExit(duration time.Duration) {
	time.Sleep(duration)
	hook.Stop()
}

func TestDiscordHook_Alert(t *testing.T) {

	err := hook.Alert("Hello World")
	if err != nil {
		panic(err)
	}
	sleepAndExit(1 * time.Second)
}

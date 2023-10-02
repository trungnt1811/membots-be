package webhook

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	discord2 "github.com/astraprotocol/affiliate-system/internal/webhook/discord"
)

// Whm is the global WebHookManager for dev-testing purposes.
var Whm *WebHookManager

// WebHookManager manages the list of WebHooks and is the main entrance to posting messages.
type WebHookManager struct {
	webHooks map[string]WebHook
}

// NewWebHookManager returns a new WebHookManager.
func NewWebHookManager(whs map[string]WebHook) *WebHookManager {
	if whs == nil {
		whs = make(map[string]WebHook)
	}
	return &WebHookManager{webHooks: whs}
}

// NewWebHookManagerFromConfig returns a new WebHookManager from the given config.
func NewWebHookManagerFromConfig(cfg conf.WebhookConfiguration) (*WebHookManager, error) {
	whs := make(map[string]WebHook)
	var err error

	discordConf := discord2.DefaultConfig()
	discordConf.Token = cfg.DcConf.Token
	discordConf.SubChannels = cfg.DcConf.SubChannels

	whs["Discord"], err = discord2.NewDiscordHook(discordConf)
	if err != nil {
		log.LG.Fatalf("failed to init discord hook: %v", err)
		return nil, err
	}

	return &WebHookManager{webHooks: whs}, nil
}

// AddWebHook adds a new WebHook to the WebHookManager.
func (m *WebHookManager) AddWebHook(name string, wh WebHook) {
	m.webHooks[name] = wh
}

// RemoveWebHook removes a WebHook given its name.
func (m *WebHookManager) RemoveWebHook(name string) {
	if m.webHooks[name] != nil {
		delete(m.webHooks, name)
	}
}

// Start starts all the WebHook of the WebHookManager.
func (m *WebHookManager) Start() error {
	var err error
	for _, wh := range m.webHooks {
		err = wh.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

// Stop terminates all the processes of the WebHookManager.
func (m *WebHookManager) Stop() {
	for _, wh := range m.webHooks {
		go wh.Stop()
	}
}

// Alert pushes the given alert message to the WebHook's.
func (m *WebHookManager) Alert(msg interface{}) error {
	for _, wh := range m.webHooks {
		_ = wh.Alert(msg)
	}

	return nil
}

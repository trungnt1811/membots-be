package discord

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// DiscordHook implements a Discord webhook for pushing messages.
type DiscordHook struct {
	*discordgo.Session
	cfg *DiscordConfig

	updateChannel chan interface{}
	closeChannel  chan interface{}
	stopped       bool

	mtx *sync.Mutex
	log zerolog.Logger
}

// NewDiscordHook creates and returns a new DiscordHook with the given DiscordConfig.
func NewDiscordHook(cfg DiscordConfig) (*DiscordHook, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %v", cfg.Token))
	if err != nil {
		return nil, err
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%v-%v]", i, "Discord"))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s", i)
		},
	}
	discordLog := zerolog.New(writer)

	return &DiscordHook{
		Session:       session,
		cfg:           &cfg,
		updateChannel: make(chan interface{}),
		closeChannel:  make(chan interface{}),
		stopped:       false,
		mtx:           new(sync.Mutex),
		log:           discordLog,
	}, nil
}

// Ping checks if the remote Tele service is alive.
func (h *DiscordHook) Ping() error {
	return nil
}

// Start starts the TeleHook.
func (h *DiscordHook) Start() error {
	logEvent := h.log.Info()
	logEvent.Timestamp().Msg("STARTED\n")
	go h.handleUpdateChan()

	return nil
}

// Stop stops the TeleHook.
func (h *DiscordHook) Stop() {
	h.closeChannel <- true

	for !h.stopped {
		time.Sleep(1 * time.Second)
	}

	// close channels
	close(h.updateChannel)
	close(h.closeChannel)

	logEvent := h.log.Info()
	logEvent.Timestamp().Msg("STOPPED!\n")
}

// Alert pushes the given alert message to the TeleHook.
func (h *DiscordHook) Alert(msg interface{}) error {
	return h.pushMessage(msg)
}

func (h *DiscordHook) pushMessage(msg interface{}) error {
	if len(h.updateChannel) >= h.cfg.MessageQueueSize {
		h.log.Error().Timestamp().Msgf("channel is full")
		return fmt.Errorf("channel is full")
	}

	h.updateChannel <- msg

	return nil
}

func (h *DiscordHook) handleUpdateChan() {
	stop := false
	for {
		select {
		case <-h.closeChannel:
			h.mtx.Lock()
			h.stopped = true
			stop = true
			h.mtx.Unlock()
			break
		case msg := <-h.updateChannel:
			h.log.Info().Timestamp().Msgf("New message received: %v\n", util.MustFormatJson(msg))
			for _, subChannel := range h.cfg.SubChannels {
				msgStr, ok := msg.(string)
				if !ok {
					msgStr = util.MustFormatJson(msg)
				}
				if _, err := h.ChannelMessageSend(subChannel, msgStr); err != nil {
					h.log.Error().Timestamp().Msgf("Error sending message to Discord.\nMessage: %v\nError: %v", msg, err)
				}
			}
		default:
			time.Sleep(1 * time.Second)
		}

		if stop {
			break
		}
	}
}

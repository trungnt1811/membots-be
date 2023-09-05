package watcher

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util"
)

const (
	BLOCKTIME = 3 * time.Second
	RETRIES   = 3
)

type Watcher struct {
	usecase *WatcherUsecase
}

func NewWatcher(watcherUsecase *WatcherUsecase) *Watcher {
	return &Watcher{watcherUsecase}
}

func (w *Watcher) Start(channel *util.Channel) {
	// go w.usecase.ListenNewBroadCastTx(channel)
}

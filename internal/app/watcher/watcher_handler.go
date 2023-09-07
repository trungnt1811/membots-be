package watcher

import (
	"github.com/astraprotocol/affiliate-system/internal/util"
	"time"
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

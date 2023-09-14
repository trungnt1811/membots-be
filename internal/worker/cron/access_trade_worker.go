package cron

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util/log"

	"github.com/go-co-op/gocron"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type AccessTradeWorker struct {
	RedSync *redsync.Redsync
	Usecase interfaces.ATUCase
}

func NewAccessTradeWorker(rdc *redis.Client, usecase interfaces.ATUCase) *AccessTradeWorker {
	pool := goredis.NewPool(rdc)
	// Create an instance of redisync to be used to obtain
	// a mutual exclusion lock
	rs := redsync.New(pool)
	return &AccessTradeWorker{
		RedSync: rs,
		Usecase: usecase,
	}
}

func (worker *AccessTradeWorker) RunJob() {
	s := gocron.NewScheduler(time.UTC)
	log.LG.Info("run accesstrade jobs")
	mutexName := "lock-accesstrade-sync"
	_, err := s.Every(2).Hours().Do(func() {
		ctx := context.Background()
		mutex := worker.RedSync.NewMutex(mutexName, redsync.WithExpiry(time.Hour))
		if err := mutex.LockContext(ctx); err != nil {
			log.LG.Infof("lock error: %v", err)
		} else {
			count, err := worker.Usecase.QueryAndSaveCampaigns(true)
			if err != nil {
				log.LG.Infof("sync campaigns error: %v", err)
			} else {
				log.LG.Infof("sync campaigns success: %d synced", count)
			}
			// After run, release the lock so other processes or threads can obtain a lock.
			if ok, err1 := mutex.UnlockContext(ctx); !ok || err1 != nil {
				log.LG.Infof("unlocked error: %v", err1)
			}
		}
	})
	if err != nil {
		log.LG.Infof("start job error: %v", err)
	}
	s.StartBlocking()
}

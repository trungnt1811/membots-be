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

type OrderListWorker struct {
	RedSync *redsync.Redsync
	UCase   interfaces.OrderUCase
}

func NewOrderListWorker(rdc *redis.Client, usecase interfaces.OrderUCase) *OrderListWorker {
	pool := goredis.NewPool(rdc)
	// Create an instance of redisync to be used to obtain
	// a mutual exclusion lock
	rs := redsync.New(pool)
	return &OrderListWorker{
		RedSync: rs,
		UCase:   usecase,
	}
}

func (worker *OrderListWorker) RunJob() {
	s := gocron.NewScheduler(time.UTC)
	log.LG.Info("run order-sync jobs")
	mutexName := "lock-order-sync"
	_, err := s.Every(4).Hours().Do(func() {
		ctx := context.Background()
		mutex := worker.RedSync.NewMutex(mutexName, redsync.WithExpiry(time.Hour))
		if err := mutex.LockContext(ctx); err != nil {
			log.LG.Infof("lock error: %v", err)
		} else {
			count, err := worker.UCase.CheckOrderListAndSync()
			if err != nil {
				log.LG.Infof("sync orders error: %v", err)
			} else {
				log.LG.Infof("sync orders success: %d synced", count)
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

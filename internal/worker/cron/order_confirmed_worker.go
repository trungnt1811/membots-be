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

type OrderConfirmedWorker struct {
	RedSync *redsync.Redsync
	UCase   interfaces.OrderUCase
}

func NewOrderConfirmedWorker(rdc *redis.Client, usecase interfaces.OrderUCase) *OrderConfirmedWorker {
	pool := goredis.NewPool(rdc)
	// Create an instance of redisync to be used to obtain
	// a mutual exclusion lock
	rs := redsync.New(pool)
	return &OrderConfirmedWorker{
		RedSync: rs,
		UCase:   usecase,
	}
}

func (worker *OrderConfirmedWorker) RunJob() {
	s := gocron.NewScheduler(time.UTC)
	log.LG.Info("run order-confirmed jobs")
	mutexName := "lock-order-confirmed-sync"
	_, err := s.Every(12).Hours().Do(func() {
		ctx := context.Background()
		mutex := worker.RedSync.NewMutex(mutexName, redsync.WithExpiry(time.Hour))
		if err := mutex.LockContext(ctx); err != nil {
			log.LG.Infof("lock error: %v", err)
		} else {
			count, err := worker.UCase.CheckOrderConfirmed()
			if err != nil {
				log.LG.Infof("check order confirmed error: %v", err)
			} else {
				log.LG.Infof("check order confirmed success: %d synced", count)
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

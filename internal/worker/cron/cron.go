package cron

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"gorm.io/gorm"
)

func RegisterCronJobs(config *conf.Configuration, db *gorm.DB) {
	rdc := conf.RedisConn()

	appNotiQueue := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_NOTI_APP_MESSAGE)

	// Initialize dependencies
	orderRepo := order.NewOrderRepository(db)
	rewardRepo := reward.NewRewardRepository(db)

	campRepo := campaign.NewCampaignRepository(db)
	atRepository := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)
	atUCase := accesstrade.NewAccessTradeUCase(atRepository, campRepo)

	// Initialize workers
	atWorker := NewAccessTradeWorker(rdc, atUCase)
	go func() {
		// Run brand's campaigns sync
		atWorker.RunJob()
	}()

	notiWorker := NewNotiScheduler(appNotiQueue, rewardRepo, orderRepo)
	notiWorker.StartNotiDailyReward()
}

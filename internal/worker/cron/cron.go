package cron

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	"gorm.io/gorm"
)

func RegisterCronJobs(config *conf.Configuration, db *gorm.DB) {
	rdc := conf.RedisConn()

	// Initialize dependencies
	campRepo := campaign.NewCampaignRepository(db)
	atRepository := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)
	atUCase := accesstrade.NewAccessTradeUCase(atRepository, campRepo)

	// Initialize workers
	atWorker := NewAccessTradeWorker(rdc, atUCase)
	go func() {
		// Run brand's campaigns sync
		atWorker.RunJob()
	}()
}

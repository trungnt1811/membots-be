package internal

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	"gorm.io/gorm"
)

// RegisterCronJobs Register internal app's cron jobs with default dependencies
func RegisterCronJobs(config *conf.Configuration, db *gorm.DB) {
	rdc := conf.RedisConn()

	// Initialize dependencies
	campRepo := campaign.NewCampaignRepository(db)
	atRepository := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)
	atUCase := accesstrade.NewAccessTradeUCase(atRepository, campRepo)

	// Initialize workers
	atWorker := accesstrade.NewAccessTradeWorker(rdc, atUCase)
	go func() {
		// Run brand's campaigns sync
		atWorker.RunJob()
	}()
}

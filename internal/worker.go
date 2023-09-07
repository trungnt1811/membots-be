package internal

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	"gorm.io/gorm"
)

// Register internal app's cron jobs with default dependencies
func RegisterCronjobs(config *conf.Configuration, db *gorm.DB) {
	rdc := conf.RedisConn()

	// Initialize dependencies
	campRepo := campaign.NewCampaignRepository(db)
	atRepository := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)
	atUsecase := accesstrade.NewAccessTradeUsecase(atRepository, campRepo)

	// Initialize workers
	atWorker := accesstrade.NewAccessTradeWorker(rdc, atUsecase)
	go func() {
		// Run brand's campaigns sync
		atWorker.RunJob()
	}()
}

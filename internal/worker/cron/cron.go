package cron

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade"
	"gorm.io/gorm"
)

func RegisterCronJobs(config *conf.Configuration, db *gorm.DB) {
	rdc := conf.RedisConn()

	// Initialize dependencies
	// campRepo := campaign.NewCampaignRepository(db)
	atRepository := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)
	// atUCase := accesstrade.NewAccessTradeUCase(atRepository, campRepo)

	orderRepo := order.NewOrderRepository(db)
	orderUCase := order.NewOrderUCase(orderRepo, atRepository)

	// Initialize workers
	// atWorker := NewAccessTradeWorker(rdc, atUCase)
	// go func() {
	// 	// Run brand's campaigns sync
	// 	atWorker.RunJob()
	// }()

	orderConfirmedWorker := NewOrderConfirmedWorker(rdc, orderUCase)
	go func() {
		// Run order confirmed sync
		orderConfirmedWorker.RunJob()
	}()
}

package kafkaconsumer

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/app/user_view_aff_camp"
	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/infra/exchange"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisConsumers(config *conf.Configuration, db *gorm.DB) {
	rdb := conf.RedisConn()
	redisClient := caching.NewRedisClient(rdb)

	// Kafka queue
	orderApproveQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_AFF_ORDER_UPDATE, msgqueue.KAFKA_GROUP_ID)
	appNotiQueue := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_NOTI_APP_MESSAGE)
	userViewAffCampQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_USER_VIEW_AFF_CAMP, msgqueue.KAFKA_GROUP_ID_USER_VIEW_AFF_CAMP)
	shippingReceiptConsumer := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_IMPORT_RECEIPT_TX, msgqueue.KAFKA_GROUP_ID)

	// Repository
	orderRepo := order.NewOrderRepository(db)
	rewardRepo := reward.NewRewardRepository(db)
	userViewAffCampRepo := user_view_aff_camp.NewUserViewAffCampRepository(db)

	tikiClient := exchange.NewTikiClient(exchange.TikiClientConfig{BaseUrl: config.Tiki.ApiUrl})
	tikiClientCache := exchange.NewTikiClientCache(tikiClient, redisClient)

	// Start consumer
	rewardMaker := NewRewardMaker(rewardRepo, orderRepo, tikiClientCache, orderApproveQueue, appNotiQueue)
	rewardMaker.ListenOrderApproved()

	userViewAffCampConsumer := NewUserViewAffCampConsumer(userViewAffCampRepo, userViewAffCampQueue)
	errChnUserViewAffCamp := userViewAffCampConsumer.StartListen()
	go func() {
		for err := range errChnUserViewAffCamp {
			log.Fatal().Err(err).Msg("failed to start userViewAffCampConsumer")
		}
	}()

	shippingReceiptListener := NewShippingReceiptListener(rewardRepo, shippingReceiptConsumer)
	shippingReceiptListener.ListenShippingReceipt()
}

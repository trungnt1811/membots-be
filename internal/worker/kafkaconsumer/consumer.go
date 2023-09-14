package kafkaconsumer

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/app/user_view_brand"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisConsumers(config *conf.Configuration, db *gorm.DB) {
	// Kafka queue
	orderApproveQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_AFF_ORDER_APPROVE, msgqueue.KAFKA_GROUP_ID)
	kafkaNotiMsgProducer := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_NOTI_APP_MESSAGE)
	userViewBrandQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_USER_VIEW_BRAND, msgqueue.KAFKA_GROUP_ID_USER_VIEW_BRAND)

	// Repository
	orderRepo := order.NewOrderRepository(db)
	rewardRepo := reward.NewRewardRepository(db)
	userViewBrandRepo := user_view_brand.NewUserViewBrandRepository(db)

	rewardMaker := NewRewardMaker(rewardRepo, orderRepo, orderApproveQueue, kafkaNotiMsgProducer)

	// Start consumer
	go rewardMaker.ListenOrderApproved()

	userViewBrandConsumer := NewUserViewBrandConsumer(userViewBrandRepo, userViewBrandQueue)
	errChnUserViewBrand := userViewBrandConsumer.StartListen()
	go func() {
		for err := range errChnUserViewBrand {
			log.Fatal().Err(err).Msg("failed to start userViewBrandConsumer")
		}
	}()
}

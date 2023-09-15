package kafkaconsumer

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/app/user_view_aff_camp"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisConsumers(config *conf.Configuration, db *gorm.DB) {
	// Kafka queue
	orderApproveQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_AFF_ORDER_APPROVE, msgqueue.KAFKA_GROUP_ID)
	kafkaNotiMsgProducer := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_NOTI_APP_MESSAGE)
	userViewAffCampQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_USER_VIEW_BRAND, msgqueue.KAFKA_GROUP_ID_USER_VIEW_BRAND)

	// Repository
	orderRepo := order.NewOrderRepository(db)
	rewardRepo := reward.NewRewardRepository(db)
	userViewAffCampRepo := user_view_aff_camp.NewUserViewAffCampRepository(db)

	rewardMaker := NewRewardMaker(rewardRepo, orderRepo, orderApproveQueue, kafkaNotiMsgProducer)

	// Start consumer
	go rewardMaker.ListenOrderApproved()

	userViewAffCampConsumer := NewUserViewAffCampConsumer(userViewAffCampRepo, userViewAffCampQueue)
	errChnUserViewAffCamp := userViewAffCampConsumer.StartListen()
	go func() {
		for err := range errChnUserViewAffCamp {
			log.Fatal().Err(err).Msg("failed to start userViewAffCampConsumer")
		}
	}()
}

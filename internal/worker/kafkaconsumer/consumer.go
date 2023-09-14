package kafkaconsumer

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"gorm.io/gorm"
)

func RegisConsumers(config *conf.Configuration, db *gorm.DB) {
	// Kafka queue
	orderApproveQueue := msgqueue.NewKafkaConsumer(msgqueue.KAFKA_TOPIC_AFF_ORDER_APPROVE, msgqueue.KAFKA_GROUP_ID)
	kafkaNotiMsgProducer := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_NOTI_APP_MESSAGE)

	// Repository
	orderRepo := order.NewOrderRepository(db)
	rewardRepo := reward.NewRewardRepository(db)

	rewardMaker := NewRewardMaker(rewardRepo, orderRepo, orderApproveQueue, kafkaNotiMsgProducer)

	// Start consumer
	go rewardMaker.ListenOrderApproved()
}

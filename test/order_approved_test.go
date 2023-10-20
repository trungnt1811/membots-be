package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestProduceOrderApprovedMsg(t *testing.T) {
	producer := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_AFF_ORDER_UPDATE)

	msg := msgqueue.MsgOrderUpdated{
		UserId:      584,
		AtOrderID:   "436263878870771Kien",
		OrderStatus: model.OrderStatusPending,
	}

	b, err := json.Marshal(msg)
	assert.Nil(t, err)

	err = producer.WriteMessages(context.Background(), kafka.Message{
		Value: b,
	})
	assert.Nil(t, err)

	fmt.Println("Push MsgOrderUpdated------------")
}

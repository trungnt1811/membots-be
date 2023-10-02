package msgqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/segmentio/kafka-go"
)

type userViewAffCampProducer struct {
	Producer *QueueWriter
	Stream   chan []*dto.UserViewAffCampDto
	stopSig  chan bool
}

func (p userViewAffCampProducer) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-p.stopSig:
				ticker.Stop()
				return
			case <-ticker.C:
				p.handler()
			}
		}
	}()
}
func (p userViewAffCampProducer) handler() {
	for {
		events := <-p.Stream
		for _, event := range events {
			b, _ := json.Marshal(event)
			msg := kafka.Message{
				Key:   []byte(fmt.Sprint(event.UserId)),
				Value: b,
			}
			_ = p.Producer.WriteMessages(context.Background(), msg)
		}

	}
}

func (p userViewAffCampProducer) Stop() {
	p.stopSig <- true
}

type Producer interface {
	Start()
	Stop()
}

func NewUserViewAffCampProducer(producer *QueueWriter, stream chan []*dto.UserViewAffCampDto) Producer {
	return &userViewAffCampProducer{
		Producer: producer,
		Stream:   stream,
		stopSig:  make(chan bool, 1),
	}
}

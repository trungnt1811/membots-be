package kafkaconsumer

import (
	"context"
	"encoding/json"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/rs/zerolog/log"
)

type UserViewAffCampConsumer struct {
	KafkaConsumer             *msgqueue.QueueReader
	UserViewAffCampRepository interfaces.UserViewAffCampRepository
}

func NewUserViewAffCampConsumer(
	repository interfaces.UserViewAffCampRepository,
	kafkaConsumer *msgqueue.QueueReader) *UserViewAffCampConsumer {
	return &UserViewAffCampConsumer{
		UserViewAffCampRepository: repository,
		KafkaConsumer:             kafkaConsumer,
	}
}

func (c *UserViewAffCampConsumer) StartListen() <-chan error {
	log.Info().Msg("start insert user view aff camp consumer")
	errChn := make(chan error)
	// Start the listener

	go func() {
		ctx := context.Background()
		for {
			m, err := c.KafkaConsumer.Reader.FetchMessage(ctx)
			if err != nil {
				log.Error().Err(err)
				continue
			}
			var payload dto.UserViewAffCampDto
			err = json.Unmarshal(m.Value, &payload)
			if err != nil {
				log.Error().Err(err)
				if err = c.KafkaConsumer.Reader.CommitMessages(ctx, m); err != nil {
					log.Error().Err(err)
				}
				continue
			}
			_ = c.UserViewAffCampRepository.CreateUserViewAffCamp(context.Background(),
				&model.UserViewAffCamp{
					UserId:    payload.UserId,
					AffCampId: payload.AffCampId,
				})

			if err = c.KafkaConsumer.Reader.CommitMessages(ctx, m); err != nil {
				log.Error().Err(err)
				continue
			}
			if err != nil {
				log.Error().Any("UserViewAffCampConsumer commit messages error", err)
				errChn <- err
			}
		}
	}()
	return errChn
}

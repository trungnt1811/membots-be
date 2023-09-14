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

type UserViewBrandConsumer struct {
	KafkaConsumer           *msgqueue.QueueReader
	UserViewBrandRepository interfaces.UserViewBrandRepository
}

func NewUserViewBrandConsumer(
	repository interfaces.UserViewBrandRepository,
	kafkaConsumer *msgqueue.QueueReader) *UserViewBrandConsumer {
	return &UserViewBrandConsumer{
		UserViewBrandRepository: repository,
		KafkaConsumer:           kafkaConsumer,
	}
}

func (c *UserViewBrandConsumer) StartListen() <-chan error {
	log.Info().Msg("start insert user view brand consumer")
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
			var payload dto.UserViewBrandDto
			err = json.Unmarshal(m.Value, &payload)
			if err != nil {
				log.Error().Err(err)
				if err = c.KafkaConsumer.Reader.CommitMessages(ctx, m); err != nil {
					log.Error().Err(err)
				}
				continue
			}
			_ = c.UserViewBrandRepository.CreateUserViewBrand(context.Background(),
				&model.UserViewBrand{
					UserId:  payload.UserId,
					BrandId: payload.BrandId,
				})

			if err = c.KafkaConsumer.Reader.CommitMessages(ctx, m); err != nil {
				log.Error().Err(err)
				continue
			}
			if err != nil {
				log.Error().Any("UserViewBrandConsumer commit messages error", err)
				errChn <- err
			}
		}
	}()
	return errChn
}

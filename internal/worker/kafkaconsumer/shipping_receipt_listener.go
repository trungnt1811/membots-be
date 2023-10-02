package kafkaconsumer

import (
	"context"
	"encoding/json"

	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
)

type ShippingReceiptListener struct {
	rewardRepo   interfaces.RewardRepository
	deliReceiptQ *msgqueue.QueueReader
}

func NewShippingReceiptListener(rewardRepo interfaces.RewardRepository,
	deliReceiptQ *msgqueue.QueueReader) ShippingReceiptListener {
	return ShippingReceiptListener{
		rewardRepo:   rewardRepo,
		deliReceiptQ: deliReceiptQ,
	}
}

func (u *ShippingReceiptListener) ListenShippingReceipt() {
	errChn := make(chan error)
	go func() {
		for err := range errChn {
			log.LG.Errorf("ShippingReceiptListener - failed to process approved order tx: %v", err)
		}
	}()

	go func() {
		log.LG.Infof("ShippingReceiptListener - Start reading new approved order")
		for {
			ctx := context.Background()
			/* ==========================================================================
			SECTION: reading message
			=========================================================================== */
			msg, err := u.deliReceiptQ.FetchMessage(ctx)
			if err != nil {
				errChn <- err
				continue
			}

			var receipt msgqueue.DeliveryMsg
			err = json.Unmarshal(msg.Value, &receipt)
			if err != nil {
				_ = u.commitDeliveryReceiptMsg(msg)
				errChn <- err
				continue
			}
			log.LG.Infof("ShippingReceiptListener - Read new shipping receipt: %v", receipt.RequestId)

			/* ==========================================================================
			SECTION: processing
			=========================================================================== */
			status := model.ShippingStatusFailed
			if receipt.TxStatus == 1 {
				status = model.ShippingStatusSuccess
			}
			err = u.rewardRepo.UpdateWithdrawShippingStatus(ctx, receipt.RequestId, receipt.TxHash, status)
			if err != nil {
				_ = u.commitDeliveryReceiptMsg(msg)
				errChn <- err
				continue
			}

			_ = u.commitDeliveryReceiptMsg(msg)
		}
	}()
}

func (u *ShippingReceiptListener) commitDeliveryReceiptMsg(message kafka.Message) error {
	err := u.deliReceiptQ.CommitMessages(context.Background(), message)
	if err != nil {
		log.LG.Errorf("Failed to commit order approved message: %v", err)
		return err
	}
	return nil
}

package kafkaconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
)

type ShippingReceiptListener struct {
	rewardRepo   interfaces.RewardRepository
	deliReceiptQ *msgqueue.QueueReader
	appNotiQ     *msgqueue.QueueWriter
}

func NewShippingReceiptListener(rewardRepo interfaces.RewardRepository,
	appNotiQ *msgqueue.QueueWriter,
	deliReceiptQ *msgqueue.QueueReader) ShippingReceiptListener {
	return ShippingReceiptListener{
		rewardRepo:   rewardRepo,
		deliReceiptQ: deliReceiptQ,
		appNotiQ:     appNotiQ,
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
			log.LG.Infof("ShippingReceiptListener - Read new shipping receipt: %v. Tx %v", receipt.RequestId, receipt.TxHash)

			/* ==========================================================================
			SECTION: processing
			=========================================================================== */
			err = u.processReceipt(ctx, &receipt)
			if err != nil {
				_ = u.commitDeliveryReceiptMsg(msg)
				errChn <- err
				continue
			}

			_ = u.commitDeliveryReceiptMsg(msg)
		}
	}()
}

func (u *ShippingReceiptListener) processReceipt(ctx context.Context, receipt *msgqueue.DeliveryMsg) error {
	status := model.ShippingStatusFailed
	if receipt.TxStatus == 1 {
		status = model.ShippingStatusSuccess
	}

	err := u.rewardRepo.UpdateWithdrawShippingStatus(ctx, receipt.RequestId, receipt.TxHash, status)
	if err != nil {
		return err
	}

	withdraw, err := u.rewardRepo.GetWithdrawByShippingRequest(ctx, receipt.RequestId)
	if err != nil {
		return err
	}

	return u.notiWithdrawSuccess(withdraw.UserId, receipt.TxHash, withdraw.Amount)
}

func (u *ShippingReceiptListener) commitDeliveryReceiptMsg(message kafka.Message) error {
	err := u.deliReceiptQ.CommitMessages(context.Background(), message)
	if err != nil {
		log.LG.Errorf("Failed to commit order approved message: %v", err)
		return err
	}
	return nil
}

func (u *ShippingReceiptListener) notiWithdrawSuccess(userId uint, txHash string, amount float64) error {
	notiAmt := util.FormatNotiAmt(amount)
	notiMsg := msgqueue.AppNotiMsg{
		Category: msgqueue.NotiCategoryWallet,
		Title:    fmt.Sprintf("Rút thành công %v ASA từ hoàn mua sắm", notiAmt),
		Body:     fmt.Sprintf("%v ASA vừa được rút về ví chính từ ví hoàn mua sắm thành công 👌🏻", notiAmt),
		UserId:   userId,
		Data:     msgqueue.GetTxDetailsNotiData(txHash),
	}

	b, err := json.Marshal(notiMsg)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(strconv.FormatUint(uint64(userId), 10)),
		Value: b,
	}

	err = u.appNotiQ.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.LG.Infof("Pushed withdraw reward noti to user %v\n", userId)

	return nil
}

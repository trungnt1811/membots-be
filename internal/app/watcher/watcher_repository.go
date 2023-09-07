package watcher

import (
	"context"
	"encoding/json"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/segmentio/kafka-go"
)

type WatcherRepository struct {
	pendingTxQ *msgqueue.MessageQueue
	TxReceiptQ *msgqueue.MessageQueue
}

func NewWatcherRepository(pendingTxQ *msgqueue.MessageQueue, txReceiptQ *msgqueue.MessageQueue) *WatcherRepository {
	return &WatcherRepository{
		pendingTxQ: pendingTxQ,
		TxReceiptQ: txReceiptQ,
	}
}

func (rp *WatcherRepository) CommitPendingTx(msg kafka.Message) error {
	return rp.pendingTxQ.Reader.CommitMessages(context.Background(), msg)
}

func (rp *WatcherRepository) PushToTxReceiptQueue(msg kafka.Message) error {
	return rp.TxReceiptQ.Writer.WriteMessages(context.Background(), msg)
}

func (rp *WatcherRepository) FetchPendingTx() (kafka.Message, util.TxInfo, error) {
	var txinfo util.TxInfo

	m, err := rp.pendingTxQ.Reader.FetchMessage(context.Background())

	if err == nil {
		err = json.Unmarshal(m.Value, &txinfo)
	}

	return m, txinfo, err
}

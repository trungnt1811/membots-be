package interfaces

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/segmentio/kafka-go"
)

type WatcherUsecase interface {
	DoCheckEvmUnconfirmTx(tx util.TxInfo) (*types.Receipt, error)
	CheckFinality(blocknumber uint64) error
	ClassifyTx(txinfo util.TxInfo, result int, res *types.Receipt, channel *util.Channel) error
	AlertFailTx(txHash common.Hash, txType int) error
	TxTimeout(txinfo util.TxInfo, result int, channel *util.Channel) error
	CheckEvmTxDelivery(txinfo util.TxInfo, message *kafka.Message, channel *util.Channel)
	ListenNewBroadCastTx(channel *util.Channel)
}

type WatcherRepository interface {
	CommitPendingTx(msg kafka.Message) error
	PushToTxReceiptQueue(msg kafka.Message) error
	FetchPendindTx() (kafka.Message, util.TxInfo, error)
}

type WatchingEvmClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

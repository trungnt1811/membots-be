package test

import (
	"context"
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/app/watcher"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

type ReceiptRs struct {
	r   *types.Receipt
	err error
}

type MockWatchingClient struct {
	currentBlock uint64
	tx           map[string]*ReceiptRs
}

func (mock *MockWatchingClient) BlockNumber(ctx context.Context) (uint64, error) {
	mock.currentBlock += 1
	return uint64(mock.currentBlock), nil
}

func (mock *MockWatchingClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return mock.tx[txHash.Hex()].r, mock.tx[txHash.Hex()].err
}

func (mock *MockWatchingClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func Test_Parse_Receipt(t *testing.T) {
	asserts := assert.New(t)

	txs := []util.TxInfo{
		{
			TxHash: common.HexToHash("abd").Hex(),
			Type:   util.TXTYPE_REWARD_SHIPPING,
		},
		{
			TxHash: common.HexToHash("abc").Hex(),
			Type:   util.TXTYPE_REWARD_SHIPPING,
		},
		{
			TxHash: common.HexToHash("def").Hex(),
			Type:   util.TXTYPE_REWARD_SHIPPING,
		},

		{
			TxHash: common.HexToHash("bcd").Hex(),
			Type:   util.TXTYPE_REWARD_SHIPPING,
		},
	}

	testdata := map[string]*ReceiptRs{
		common.HexToHash("abd").Hex(): &ReceiptRs{
			r: &types.Receipt{
				TxHash:      common.HexToHash("abd"),
				Status:      1,
				BlockNumber: big.NewInt(0),
			},
			err: fmt.Errorf("not found"),
		},
		common.HexToHash("abc").Hex(): &ReceiptRs{
			r: &types.Receipt{
				TxHash:      common.HexToHash("abc"),
				Status:      0,
				BlockNumber: big.NewInt(0),
			},
			err: nil,
		},
		common.HexToHash("def").Hex(): &ReceiptRs{
			r: &types.Receipt{
				TxHash:      common.HexToHash("def"),
				Status:      1,
				BlockNumber: big.NewInt(0),
			},
			err: nil,
		},

		common.HexToHash("bcd").Hex(): &ReceiptRs{
			r: &types.Receipt{
				TxHash:      common.HexToHash("bcd"),
				Status:      1,
				BlockNumber: big.NewInt(100),
			},
			err: nil,
		},
	}

	mockEvmClient := MockWatchingClient{
		currentBlock: 1,
		tx:           testdata,
	}

	pendingTxQueue := msgqueue.NewMsgQueue(msgqueue.KAFKA_TOPIC_PEDNING_TX, msgqueue.KAFKA_GROUP_ID)
	txReceiptQueue := msgqueue.NewMsgQueue(msgqueue.KAFKA_TOPIC_IMPORT_RECEIPT_TX, msgqueue.KAFKA_GROUP_ID)

	watcherRepo := watcher.NewWatcherRepository(pendingTxQueue, txReceiptQueue)
	watcherUsecase := watcher.NewWatcherUsecase(&mockEvmClient, watcherRepo)
	// w := watcher.NewWatcher(watcherUsecase)

	channel := util.NewChannel()

	// Check time out case
	for _, tx := range txs {
		log.Printf("Tx %v", tx)
		watcherUsecase.CheckEvmTxDelivery(tx, nil, channel)

		rc := <-channel.TxReceipt
		log.Printf("Hash: %v Receipt: code %v, status: %v", rc.Receipt.TxHash.Hex(), rc.Code, rc.Receipt.Status)
		if rc.Receipt.TxHash == common.HexToHash("abc") {
			asserts.Equal(rc.Code, util.SUCCESS)
		}

		if rc.Receipt.TxHash == common.HexToHash("def") {
			asserts.Equal(rc.Code, util.SUCCESS)
			asserts.Equal(rc.Receipt.Status, uint64(1))
		}

		if rc.Receipt.TxHash == common.HexToHash("abd") {
			asserts.Equal(rc.Code, util.TIMEOUT)
		}

		if rc.Receipt.TxHash == common.HexToHash("bcd") {
			asserts.Equal(rc.Code, util.TIMEOUT)
		}
	}

}

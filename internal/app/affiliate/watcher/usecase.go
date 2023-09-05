package watcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/astraprotocol/affiliate-system/conf"

	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util/log"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/webhook"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/segmentio/kafka-go"
)

type WatcherUsecase struct {
	evmclient interfaces.WatchingEvmClient
	repo      interfaces.WatcherRepository
}

func NewWatcherUsecase(evmclient interfaces.WatchingEvmClient, repo interfaces.WatcherRepository) *WatcherUsecase {
	return &WatcherUsecase{
		evmclient: evmclient,
		repo:      repo,
	}
}

func (w *WatcherUsecase) ListenNewBroadCastTx(channel *util.Channel) {
	for {
		m, txinfo, err := w.repo.FetchPendindTx()
		if err != nil {
			log.LG.Errorf("pendingTxQ.Reader.FetchMessage error: %v", err.Error())
			continue
		}

		w.CheckEvmTxDelivery(txinfo, &m, channel)

		err = w.repo.CommitPendingTx(m)
		if err != nil {
			log.LG.Errorf("CommitMessages in pendingTxQ Err: %v", err.Error())
		}
	}
}

func (w *WatcherUsecase) CheckEvmTxDelivery(txinfo util.TxInfo, message *kafka.Message, channel *util.Channel) {
	conf := conf.GetConfiguration()

	// Ticker
	ticker := time.NewTicker(2 * BLOCKTIME)

	go func() {
		for range ticker.C {
			res, err := w.DoCheckEvmUnconfirmTx(txinfo)

			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					// Tx is not included in block
					continue
				} else {
					// Tx error
					log.LG.Errorf("Query receipt failed with err: %v", err.Error())
					return
				}
			} else if res != nil {
				// Tx included in some block
				blocknumner := *res.BlockNumber
				isFinalized, err := w.checkFinality(blocknumner.Uint64())
				if err != nil {
					// Tx error
					log.LG.Errorf("failed to check if block %v finalized err: %v", blocknumner.Uint64(), err.Error())
					return
				}

				if isFinalized {
					w.ClassifyTx(util.SUCCESS, res, txinfo, message, channel)
				} else {
					w.TxTimeout(util.TIMEOUT, txinfo, message, channel)
				}

				// Alert if tx fail
				if res.Status == util.FAILED {
					if conf.Env != "prod" {
						err := w.AlertFailTx(res.TxHash, txinfo.Type)
						if err != nil {
							log.LG.Errorf("Failed to alert failed tx %s to discord: %v", res.TxHash.Hex(), err)
						}
					}
					log.LG.Errorf("Transaction failed %s", res.TxHash.Hex(), err)
				}

				return
			}
		}
	}()
}

func (w *WatcherUsecase) DoCheckEvmUnconfirmTx(tx util.TxInfo) (*types.Receipt, error) {
	txhash := common.HexToHash(tx.TxHash)

	if tx.Type == util.TXTYPE_REWARD_SHIPPING {
		log.LG.Infof("Shipping reward txhash: %v", txhash.String())
	} else if tx.Type == util.TXTYPE_SELLER_WITHDRAW {
		log.LG.Infof("Seller withdraw txhash: %v", txhash.String())
	}

	return w.evmclient.TransactionReceipt(context.Background(), txhash)
}

func (w *WatcherUsecase) AlertFailTx(txHash common.Hash, txType int) error {
	discordMsg := webhook.TxFailMsg{
		TxHash:   txHash.Hex(),
		TxAction: "Reward",
	}
	if txType == util.TXTYPE_SELLER_WITHDRAW {
		discordMsg.TxAction = "WithdrawExpiredReward"
	}

	err := webhook.Whm.Alert(discordMsg.String())
	if err != nil {
		return fmt.Errorf("failed to send tx fail alert to webhook: %v", err)
	}

	return nil
}

func (w *WatcherUsecase) checkFinality(blocknumber uint64) (bool, error) {
	count := 0
	ticker := time.NewTicker(BLOCKTIME)

	for range ticker.C {
		latestblocknumber, err := w.evmclient.BlockNumber(context.Background())
		if err != nil {
			return false, err
		}

		if latestblocknumber > blocknumber {
			// finality
			return true, nil
		} else {
			count++
			if count == 10 {
				break
			}
		}

	}

	return false, errors.New("time out")
}

func (w *WatcherUsecase) TxTimeout(result int, txinfo util.TxInfo, message *kafka.Message, channel *util.Channel) error {

	switch txinfo.Type {

	case util.TXTYPE_REWARD_SHIPPING:
		res := &types.Receipt{
			TxHash: common.HexToHash(txinfo.TxHash),
			Status: util.FAILED,
		}

		w.PushToTxReceiptQueue(
			&util.TxReceiptSt{
				Code:    result,
				Receipt: res,
				Message: message,
				Hodlers: txinfo.Holders})

	case util.TXTYPE_SELLER_WITHDRAW:
		channel.Seller <- &util.TxReceiptSt{
			Code:    result,
			Receipt: &types.Receipt{TxHash: common.HexToHash(txinfo.TxHash)},
			Message: message,
			Hodlers: txinfo.Holders}
	}

	return nil
}

func (w *WatcherUsecase) ClassifyTx(result int, res *types.Receipt, txinfo util.TxInfo, message *kafka.Message, channel *util.Channel) error {

	switch txinfo.Type {

	case util.TXTYPE_REWARD_SHIPPING:
		w.PushToTxReceiptQueue(
			&util.TxReceiptSt{
				Code:    result,
				Receipt: res,
				Message: message,
				Hodlers: txinfo.Holders})

	case util.TXTYPE_SELLER_WITHDRAW:
		channel.Seller <- &util.TxReceiptSt{
			Code:    result,
			Receipt: res,
			Message: message,
			Hodlers: txinfo.Holders}
	}

	return nil
}

func (w *WatcherUsecase) PushToTxReceiptQueue(r *util.TxReceiptSt) error {
	txhash := r.Receipt.TxHash.String()

	// push Deli_Receipts to Kafka queue
	var deli_msg = util.DeliveryMsg{
		TxHash:   txhash,
		TxStatus: r.Receipt.Status,
	}

	b, err := json.Marshal(deli_msg)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: b,
	}

	return w.repo.PushToTxReceiptQueue(msg)

}

func (w *WatcherUsecase) doCheckMultiAccBalanceAndAlert(couponOperatorsAddr []string, rewardOperatorsAddr []string, deployerAddr string, minThreshold float64) {
	w.checkAccBalanceAndAlert(couponOperatorsAddr, minThreshold, util.COUPON_OPERATORS_ADDR)
	w.checkAccBalanceAndAlert([]string{deployerAddr}, minThreshold, util.DEPLOYER_ADDR)
	w.checkAccBalanceAndAlert(rewardOperatorsAddr, minThreshold, util.REWARD_OPERATORS_ADDR)
}

func (w *WatcherUsecase) checkAccBalanceAndAlert(address []string, thresholdASA float64, addr_type string) error {
	for _, addr := range address {
		low, balance, err := util.IsLowBalance(w.evmclient.(*ethclient.Client), common.HexToAddress(addr), thresholdASA)
		if err != nil {
			log.LG.Errorf("failed to check low balance for account %v: %v", addr, err)
			continue
		}

		if low {
			discordMsg := webhook.LowBalanceMsg{
				AddrType:   addr_type,
				Address:    addr,
				Balance:    balance,
				MinRequire: thresholdASA,
			}

			err = webhook.Whm.Alert(discordMsg.String())
			if err != nil {
				return fmt.Errorf("failed to send low balance alert to webhook: %v", err)
			}

		}
	}

	return nil
}

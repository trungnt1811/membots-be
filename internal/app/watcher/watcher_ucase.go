package watcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util2 "github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	webhook2 "github.com/astraprotocol/affiliate-system/internal/webhook"
	"strings"

	"github.com/astraprotocol/affiliate-system/conf"

	"time"

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

func (w *WatcherUsecase) ListenNewBroadCastTx(channel *util2.Channel) {
	for {
		m, txinfo, err := w.repo.FetchPendingTx()
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

func (w *WatcherUsecase) CheckEvmTxDelivery(txinfo util2.TxInfo, message *kafka.Message, channel *util2.Channel) {
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
					w.ClassifyTx(util2.SUCCESS, res, txinfo, message, channel)
				} else {
					w.TxTimeout(util2.TIMEOUT, txinfo, message, channel)
				}

				// Alert if tx fail
				if res.Status == util2.FAILED {
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

func (w *WatcherUsecase) DoCheckEvmUnconfirmTx(tx util2.TxInfo) (*types.Receipt, error) {
	txhash := common.HexToHash(tx.TxHash)

	if tx.Type == util2.TXTYPE_REWARD_SHIPPING {
		log.LG.Infof("Shipping reward txhash: %v", txhash.String())
	} else if tx.Type == util2.TXTYPE_SELLER_WITHDRAW {
		log.LG.Infof("Seller withdraw txhash: %v", txhash.String())
	}

	return w.evmclient.TransactionReceipt(context.Background(), txhash)
}

func (w *WatcherUsecase) AlertFailTx(txHash common.Hash, txType int) error {
	discordMsg := webhook2.TxFailMsg{
		TxHash:   txHash.Hex(),
		TxAction: "Reward",
	}
	if txType == util2.TXTYPE_SELLER_WITHDRAW {
		discordMsg.TxAction = "WithdrawExpiredReward"
	}

	err := webhook2.Whm.Alert(discordMsg.String())
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

func (w *WatcherUsecase) TxTimeout(result int, txinfo util2.TxInfo, message *kafka.Message, channel *util2.Channel) error {

	switch txinfo.Type {

	case util2.TXTYPE_REWARD_SHIPPING:
		res := &types.Receipt{
			TxHash: common.HexToHash(txinfo.TxHash),
			Status: util2.FAILED,
		}

		w.PushToTxReceiptQueue(
			&util2.TxReceiptSt{
				Code:    result,
				Receipt: res,
				Message: message,
				Hodlers: txinfo.Holders})

	case util2.TXTYPE_SELLER_WITHDRAW:
		channel.Seller <- &util2.TxReceiptSt{
			Code:    result,
			Receipt: &types.Receipt{TxHash: common.HexToHash(txinfo.TxHash)},
			Message: message,
			Hodlers: txinfo.Holders}
	}

	return nil
}

func (w *WatcherUsecase) ClassifyTx(result int, res *types.Receipt, txinfo util2.TxInfo, message *kafka.Message, channel *util2.Channel) error {

	switch txinfo.Type {

	case util2.TXTYPE_REWARD_SHIPPING:
		w.PushToTxReceiptQueue(
			&util2.TxReceiptSt{
				Code:    result,
				Receipt: res,
				Message: message,
				Hodlers: txinfo.Holders})

	case util2.TXTYPE_SELLER_WITHDRAW:
		channel.Seller <- &util2.TxReceiptSt{
			Code:    result,
			Receipt: res,
			Message: message,
			Hodlers: txinfo.Holders}
	}

	return nil
}

func (w *WatcherUsecase) PushToTxReceiptQueue(r *util2.TxReceiptSt) error {
	txhash := r.Receipt.TxHash.String()

	// push Deli_Receipts to Kafka queue
	var deli_msg = util2.DeliveryMsg{
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
	w.checkAccBalanceAndAlert(couponOperatorsAddr, minThreshold, util2.COUPON_OPERATORS_ADDR)
	w.checkAccBalanceAndAlert([]string{deployerAddr}, minThreshold, util2.DEPLOYER_ADDR)
	w.checkAccBalanceAndAlert(rewardOperatorsAddr, minThreshold, util2.REWARD_OPERATORS_ADDR)
}

func (w *WatcherUsecase) checkAccBalanceAndAlert(address []string, thresholdASA float64, addr_type string) error {
	for _, addr := range address {
		low, balance, err := util2.IsLowBalance(w.evmclient.(*ethclient.Client), common.HexToAddress(addr), thresholdASA)
		if err != nil {
			log.LG.Errorf("failed to check low balance for account %v: %v", addr, err)
			continue
		}

		if low {
			discordMsg := webhook2.LowBalanceMsg{
				AddrType:   addr_type,
				Address:    addr,
				Balance:    balance,
				MinRequire: thresholdASA,
			}

			err = webhook2.Whm.Alert(discordMsg.String())
			if err != nil {
				return fmt.Errorf("failed to send low balance alert to webhook: %v", err)
			}

		}
	}

	return nil
}

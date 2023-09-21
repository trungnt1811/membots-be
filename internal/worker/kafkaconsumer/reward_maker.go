package kafkaconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
)

type RewardMaker struct {
	rewardRepo interfaces.RewardRepository
	orderRepo  interfaces.OrderRepository
	priceRepo  interfaces.TokenPriceRepo
	approveQ   *msgqueue.QueueReader
	appNotiQ   *msgqueue.QueueWriter
}

func NewRewardMaker(rewardRepo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository,
	priceRepo interfaces.TokenPriceRepo,
	approveQ *msgqueue.QueueReader,
	appNotiQ *msgqueue.QueueWriter) *RewardMaker {
	return &RewardMaker{
		rewardRepo: rewardRepo,
		orderRepo:  orderRepo,
		priceRepo:  priceRepo,
		approveQ:   approveQ,
		appNotiQ:   appNotiQ,
	}
}

func (u *RewardMaker) ListenOrderApproved() {
	errChn := make(chan error)
	go func() {
		for err := range errChn {
			log.LG.Errorf("Reward Maker - failed to process approved order tx: %v", err)
		}
	}()

	go func() {
		log.LG.Infof("Reward Maker - Start reading new approved order")
		for {
			ctx := context.Background()
			/* ==========================================================================
			SECTION: reading message
			=========================================================================== */
			msg, err := u.approveQ.FetchMessage(ctx)
			if err != nil {
				errChn <- err
				continue
			}

			var orderApprovedMsg msgqueue.MsgOrderApproved
			err = json.Unmarshal(msg.Value, &orderApprovedMsg)
			if err != nil {
				_ = u.commitOrderApprovedMsg(msg)
				errChn <- err
				continue
			}
			newAtOrderId := orderApprovedMsg.AtOrderID
			log.LG.Infof("Read new order approved: %v", newAtOrderId)

			/* ==========================================================================
			SECTION: processing
			=========================================================================== */
			order, err := u.orderRepo.FindOrderByAccessTradeId(newAtOrderId)
			if err != nil {
				_ = u.commitOrderApprovedMsg(msg)
				errChn <- err
				continue
			}
			if order.OrderStatus != model.OrderStatusApproved {
				continue
			}

			rewardAmount, err := u.CalculateRewardAmt(float64(order.PubCommission), reward.AffCommissionFee)
			if err != nil {
				_ = u.commitOrderApprovedMsg(msg)
				errChn <- err
				continue
			}

			now := time.Now()
			newReward := model.Reward{
				UserId:         order.UserId,
				AtOrderID:      newAtOrderId,
				Amount:         rewardAmount,
				RewardedAmount: 0,
				CommissionFee:  reward.AffCommissionFee,
				EndedAt:        now.Add(reward.RewardLockTime * time.Hour),
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			err = u.rewardRepo.CreateReward(ctx, &newReward)
			if err != nil {
				_ = u.commitOrderApprovedMsg(msg)
				errChn <- err
				continue
			}

			err = u.notiOrderApproved(order.UserId, newAtOrderId, order.Merchant)
			if err != nil {
				_ = u.commitOrderApprovedMsg(msg)
				errChn <- err
				continue
			}

			_ = u.commitOrderApprovedMsg(msg)
		}
	}()
}

func (u *RewardMaker) notiOrderApproved(userId uint, atOrderId string, merchant string) error {
	notiMsg := msgqueue.AppNotiMsg{
		Category: msgqueue.NotiCategoryCommerce,
		Title:    fmt.Sprintf("Đơn hàng từ %v đã được xác nhận thành công", merchant),
		Body:     fmt.Sprintf("Đơn hàng %v đã được xác nhận thành công, bạn có thể nhận thường ngay bây giờ", atOrderId),
		UserId:   userId,
		Data:     msgqueue.GetOrderApprovedNotiData(),
	}

	b, err := json.Marshal(notiMsg)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(strconv.FormatUint(uint64(userId), 10)),
		Value: b,
	}

	err = u.appNotiQ.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.LG.Infof("Pushed Order Approved Msg to queue %v\n", notiMsg)

	return nil
}

func (u *RewardMaker) CalculateRewardAmt(affCommission float64, commissionFee float64) (float64, error) {
	astraPrice, err := u.priceRepo.GetAstraPrice(context.Background())
	if err != nil {
		return 0, err
	}
	tokenCommission := affCommission / float64(astraPrice) * (100 - commissionFee) / 100
	return util.RoundFloat(tokenCommission, 2), nil
}

func (u *RewardMaker) commitOrderApprovedMsg(message kafka.Message) error {
	err := u.approveQ.CommitMessages(context.Background(), message)
	if err != nil {
		log.LG.Errorf("Failed to commit order approved message: %v", err)
		return err
	}
	return nil
}

package kafkaconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/segmentio/kafka-go"
)

type RewardMaker struct {
	rewardRepo interfaces.RewardRepository
	orderRepo  interfaces.OrderRepository
	approveQ   *msgqueue.QueueReader
	appNotiQ   *msgqueue.QueueWriter
}

func NewRewardMaker(rewardRepo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository,
	approveQ *msgqueue.QueueReader,
	appNotiQ *msgqueue.QueueWriter) *RewardMaker {
	return &RewardMaker{
		rewardRepo: rewardRepo,
		orderRepo:  orderRepo,
		approveQ:   approveQ,
		appNotiQ:   appNotiQ,
	}
}

func (u *RewardMaker) ListenOrderApproved() {
	for {
		ctx := context.Background()

		newAtOrderId, err := u.getNewApprovedAtOrderId(ctx)
		if err != nil {
			log.Println("Failed to comsume new AccessTrade order id", err)
			continue
		}

		order, err := u.orderRepo.FindOrderByAccessTradeId(newAtOrderId)
		if err != nil {
			log.Println("Failed to get affiliate order", err)
			continue
		}

		now := time.Now()
		newReward := model.Reward{
			UserId:         order.UserId,
			AtOrderID:      newAtOrderId,
			Amount:         float64(order.PubCommission) * reward.AffCommissionFee / 100,
			RewardedAmount: 0,
			CommissionFee:  reward.AffCommissionFee,
			EndedAt:        now.Add(reward.RewardLockTime * time.Hour),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		err = u.rewardRepo.CreateReward(ctx, &newReward)
		if err != nil {
			log.Println("Failed to get affiliate order", err)
			continue
		}

		err = u.notiOrderApproved(order.UserId, newAtOrderId, order.Merchant)
		if err != nil {
			log.Println("Failed to send reward approved noti", err)
			continue
		}
	}
}

func (u *RewardMaker) getNewApprovedAtOrderId(ctx context.Context) (string, error) {
	var msg msgqueue.MsgOrderApproved
	m, err := u.approveQ.FetchMessage(ctx)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(m.Value, &msg)
	if err != nil {
		return "", err
	}

	return msg.AtOrderID, nil
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

	log.Printf("Pushed Order Approved Msg to queue %v\n", notiMsg)

	return nil
}

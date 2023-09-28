package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	scheduler "github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
)

const (
	DayToCouponExpire   = 5
	DailyRewardNotiTime = "45 09 * * *" // every day at 9:45
)

type NotiScheduler struct {
	scheduler  *scheduler.Cron
	appNotiQ   *msgqueue.QueueWriter
	rewardRepo interfaces.RewardRepository
	orderRepo  interfaces.OrderRepository
}

func NewCouponNoti(appNotiQ *msgqueue.QueueWriter,
	rewardRepo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository) *NotiScheduler {
	parser := scheduler.NewParser(
		scheduler.SecondOptional | scheduler.Minute | scheduler.Hour | scheduler.Dom | scheduler.Month | scheduler.Dow | scheduler.Descriptor,
	)
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	couponNoti := &NotiScheduler{
		scheduler:  scheduler.New(scheduler.WithParser(parser), scheduler.WithLocation(location)),
		appNotiQ:   appNotiQ,
		rewardRepo: rewardRepo,
		orderRepo:  orderRepo,
	}

	couponNoti.scheduler.Start()

	return couponNoti
}

func (n *NotiScheduler) StartNotiExpireCoupons() {
	_, err := n.scheduler.AddFunc(DailyRewardNotiTime, n.notiRewardInDay())
	if err != nil {
		log.LG.Errorf("Failed to schedule notiRewardInDay. Err %v", err)
	}
}

func (n *NotiScheduler) notiRewardInDay() func() {
	return func() {
		log.LG.Info("Start to screening user have aff reward")
		ctx := context.Background()

		// TODO: find users

		users := []uint32{584}
		for _, userId := range users {
			rewards, err := n.rewardRepo.GetInProgressRewards(ctx, userId)
			if err != nil {
				log.LG.Errorf("Failed to GetInProgressRewards of user %v. Err %v", userId, err)
				continue
			}

			var totalRewardInDay float64 = 0
			var orderCount uint = 0
			for _, r := range rewards {
				now := time.Now()
				if now.Before(r.EndAt) && now.After(r.StartAt.Add(model.OneDay)) {
					orderCount++
					totalRewardInDay += r.OneDayReward()
				}
			}

			err = n.pushDailyRewardNotiToQueue(userId, totalRewardInDay, orderCount)
			if err != nil {
				log.LG.Errorf("Failed to pushDailyRewardNotiToQueue for user %v. Err %v", userId, err)
			}
		}

	}
}

func (n *NotiScheduler) pushDailyRewardNotiToQueue(userId uint32, amount float64, orderCount uint) error {
	msg := msgqueue.AppNotiMsg{
		Title:    fmt.Sprintf("Nhận %v ASA từ hoàn mua sắm", amount),
		Body:     fmt.Sprintf("Bạn nhận được %v ASA từ %v đơn hàng hoàn mua sắm hôm nay 🤗", amount, orderCount),
		UserId:   uint(userId),
		Category: msgqueue.NotiCategoryAffiliate,
		Data:     msgqueue.GetDailyRewardNotiData(),
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	err = n.appNotiQ.WriteMessages(context.Background(), kafka.Message{
		Value: b,
	})
	if err != nil {
		return err
	}

	log.LG.Infof("Pushed coupon expire noti to queue. User %v", userId)

	return nil
}

package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
)

const (
	DailyRewardNotiTime     = "45 09 * * *" // every day at 9:45
	UserProcessingBatchSize = 200
)

type NotiScheduler struct {
	cron       *cron.Cron
	appNotiQ   *msgqueue.QueueWriter
	rewardRepo interfaces.RewardRepository
	orderRepo  interfaces.OrderRepository
}

func NewNotiScheduler(appNotiQ *msgqueue.QueueWriter,
	rewardRepo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository) *NotiScheduler {
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	// Works only if this host supports timezone
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	notiScheduler := &NotiScheduler{
		cron:       cron.New(cron.WithParser(parser), cron.WithLocation(location)),
		appNotiQ:   appNotiQ,
		rewardRepo: rewardRepo,
		orderRepo:  orderRepo,
	}

	notiScheduler.cron.Start()

	return notiScheduler
}

func (n *NotiScheduler) StartNotiDailyReward() {
	_, err := n.cron.AddFunc(DailyRewardNotiTime, n.notiRewardInDay())
	if err != nil {
		log.LG.Errorf("Failed to schedule notiRewardInDay. Err %v", err)
	}
}

func (n *NotiScheduler) notiRewardInDay() func() {
	return func() {
		log.LG.Info("Start to screening user have aff reward")
		ctx := context.Background()

		// find all users have 'rewarding' order
		users, err := n.rewardRepo.GetUsersHaveInProgressRewards(ctx)
		if err != nil {
			log.LG.Errorf("Failed to GetUsersHaveInProgressRewards. Err %v", err)
			return
		}

		// batching them
		userBatches := [][]uint32{}
		batch := []uint32{}
		for idx, userId := range users {
			batch = append(batch, userId)
			if len(batch) == UserProcessingBatchSize || idx == len(users)-1 {
				userBatches = append(userBatches, batch)
				batch = []uint32{} // we reassign 'batch', not change the 'batch' member. So it's not affect 'userBatches' member
			}
		}

		for idx, userIds := range userBatches {
			fmt.Println("Batch ", idx, ":", userIds)
			rewards, err := n.rewardRepo.GetInProgressRewardsOfMultipleUsers(ctx, userIds)
			if err != nil {
				log.LG.Errorf("Failed to GetInProgressRewardsOfMultipleUsers. Err %v", err)
				continue
			}

			rewardsByUser := map[uint32][]model.Reward{}
			for _, item := range rewards {
				rewardsByUser[uint32(item.UserId)] = append(rewardsByUser[uint32(item.UserId)], item)
			}

			// Processing each user in a single batch
			notiMessages := []kafka.Message{}
			for userId, userRewards := range rewardsByUser {
				fmt.Println("PENDING REWARDS", len(userRewards))
				data, _ := json.Marshal(userRewards)
				fmt.Println("PENDING REWARDS", string(data))

				var totalRewardInDay float64 = 0
				var orderCount uint = 0
				for _, r := range userRewards {
					now := time.Now()
					if now.Before(r.EndAt) && now.After(r.StartAt.Add(model.OneDay)) {
						orderCount++
						totalRewardInDay += r.OneDayReward()
					}
				}

				if totalRewardInDay > 0 {
					notiAmount := util.RoundFloat(totalRewardInDay, 2)
					msg := msgqueue.AppNotiMsg{
						Title:    fmt.Sprintf("Nhận %v ASA từ hoàn mua sắm", notiAmount),
						Body:     fmt.Sprintf("Bạn nhận được %v ASA từ %v đơn hàng hoàn mua sắm hôm nay 🤗", notiAmount, orderCount),
						UserId:   uint(userId),
						Category: msgqueue.NotiCategoryAffiliate,
						Data:     msgqueue.GetDailyRewardNotiData(),
					}
					b, err := json.Marshal(&msg)
					if err != nil {
						log.LG.Errorf("Failed to marshal msg %v. Err %v", msg, err)
						continue
					}
					notiMessages = append(notiMessages, kafka.Message{Value: b})
				}
			}

			if len(notiMessages) > 0 {
				err = n.pushDailyRewardNotiToQueue(notiMessages)
				if err != nil {
					log.LG.Errorf("Failed to pushDailyRewardNotiToQueue. Err %v", err)
				}
			}
		}

	}
}

func (n *NotiScheduler) pushDailyRewardNotiToQueue(msgs []kafka.Message) error {
	err := n.appNotiQ.WriteMessages(context.Background(), msgs...)
	if err != nil {
		return err
	}

	log.LG.Infof("Pushed daily reward noti to queue")

	return nil
}

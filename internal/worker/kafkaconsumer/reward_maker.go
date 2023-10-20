package kafkaconsumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/infra/alert"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/segmentio/kafka-go"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type RewardMaker struct {
	rewardRepo   interfaces.RewardRepository
	orderRepo    interfaces.OrderRepository
	priceRepo    interfaces.TokenPriceRepo
	orderUpdateQ *msgqueue.QueueReader
	appNotiQ     *msgqueue.QueueWriter
	alertClient  *alert.AlertClient
}

func NewRewardMaker(rewardRepo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository,
	priceRepo interfaces.TokenPriceRepo,
	orderUpdateQ *msgqueue.QueueReader,
	appNotiQ *msgqueue.QueueWriter,
	alertClient *alert.AlertClient) *RewardMaker {
	return &RewardMaker{
		rewardRepo:   rewardRepo,
		orderRepo:    orderRepo,
		priceRepo:    priceRepo,
		orderUpdateQ: orderUpdateQ,
		appNotiQ:     appNotiQ,
		alertClient:  alertClient,
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
		log.LG.Infof("Reward Maker - Start reading order update")
		for {
			ctx := context.Background()
			/* ==========================================================================
			SECTION: reading message
			=========================================================================== */
			msg, err := u.orderUpdateQ.FetchMessage(ctx)
			if err != nil {
				errChn <- err
				continue
			}

			var orderApprovedMsg msgqueue.MsgOrderUpdated
			err = json.Unmarshal(msg.Value, &orderApprovedMsg)
			if err != nil {
				_ = u.commitOrderUpdateMsg(msg)
				errChn <- err
				continue
			}
			log.LG.Infof("Read new order updated: %v", orderApprovedMsg.AtOrderID)

			/* ==========================================================================
			SECTION: processing
			=========================================================================== */
			err = u.processOrderUpdateMsg(ctx, orderApprovedMsg)
			if err != nil {
				_ = u.commitOrderUpdateMsg(msg)
				errChn <- err
				continue
			}

			_ = u.commitOrderUpdateMsg(msg)
		}
	}()
}

func (u *RewardMaker) processOrderUpdateMsg(ctx context.Context, msg msgqueue.MsgOrderUpdated) error {
	stellaCommission := conf.GetConfiguration().Aff.StellaCommission

	newAtOrderId := msg.AtOrderID
	order, err := u.orderRepo.FindOrderByAccessTradeId(newAtOrderId)
	if err != nil {
		return err
	}

	rewardAmount, err := u.CalculateRewardAmt(float64(order.PubCommission), stellaCommission)
	if err != nil {
		return err
	}

	var BeforeRewardingStatuses = []string{model.OrderStatusInitial, model.OrderStatusPending, model.OrderStatusApproved}
	if slices.Contains(BeforeRewardingStatuses, order.OrderStatus) {
		currentRw, err := u.rewardRepo.GetRewardByAtOrderId(ctx, newAtOrderId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			now := time.Now()
			newReward := model.Reward{
				UserId:           order.UserId,
				AtOrderID:        newAtOrderId,
				Amount:           rewardAmount,
				RewardedAmount:   0,
				CommissionFee:    stellaCommission,
				ImmediateRelease: model.ImmediateRelease,
				StartAt:          now,
				EndAt:            now.Add(reward.RewardLockTime * time.Hour),
			}
			err = u.rewardRepo.CreateReward(ctx, &newReward)
			if err != nil {
				return err
			}

		} else {
			if rewardAmount != currentRw.Amount {
				log.LG.Infof("Order %v updated reward amount", newAtOrderId)
				now := time.Now()
				rewardUpdates := &model.Reward{
					Amount:  rewardAmount,
					StartAt: now,
					EndAt:   now.Add(reward.RewardLockTime * time.Hour),
				}
				err = u.rewardRepo.UpdateRewardByAtOrderId(newAtOrderId, rewardUpdates)
				if err != nil {
					return err
				}
			}
		}
	}

	if order.OrderStatus == model.OrderStatusApproved {
		_, err = u.orderRepo.UpdateOrder(&model.AffOrder{ID: order.ID, OrderStatus: model.OrderStatusRewarding})
		if err != nil {
			return err
		}
	}

	if msg.OrderStatus != order.OrderStatus {
		return u.notiOrderStatus(order.UserId, order.ID, order.OrderStatus, newAtOrderId, order.Merchant, rewardAmount)
	}

	if rewardAmount >= reward.SuspiciousRewardAmount {
		alertMsg := alert.SuspiciousOrderMsg{
			OrderId:   order.AccessTradeOrderId,
			Billing:   uint(order.Billing),
			Reward:    rewardAmount,
			Threshold: reward.SuspiciousRewardAmount,
		}
		log.LG.Errorf("suspicious order: %v", alertMsg.String())

		err := u.alertClient.SendMessage(alertMsg.String())
		if err != nil {
			log.LG.Errorf("failed to alert suspicious order: %v. Err: %v", order.AccessTradeOrderId, err)
		}
	}

	return nil
}

func (u *RewardMaker) notiOrderStatus(userId uint, orderId uint, orderStatus, atOrderId, merchant string, rewardAmount float64) error {
	title := ""
	body := ""
	notiAmt := util.FormatNotiAmt(rewardAmount)

	switch orderStatus {
	case model.OrderStatusInitial, model.OrderStatusPending:
		title = fmt.Sprintf("Đơn hoàn mua sắm mới từ %v", merchant)
		body = fmt.Sprintf("Đơn hàng #%v của bạn vừa được cập nhật. Bấm để xem chi tiết!", atOrderId)
	case model.OrderStatusApproved:
		title = "Đơn hoàn mua sắm được xác nhận"
		body = fmt.Sprintf("%v ASA sẽ được hoàn cho đơn %v #%v vừa xác nhận hoàn tất 😝", notiAmt, merchant, atOrderId)
	case model.OrderStatusCancelled:
		title = "Đơn hoàn mua sắm đã huỷ"
		body = fmt.Sprintf("Đơn hàng %v #%v của bạn đã huỷ. Bấm để xem chi tiết!", merchant, atOrderId)
	case model.OrderStatusRejected:
		title = "Đơn hoàn mua sắm bị từ chối"
		body = fmt.Sprintf("Đơn hàng %v #%v của bạn đã bị từ chối hoàn ASA từ đối tác. Bấm để xem chi tiết!", merchant, atOrderId)
	default:
		// model.OrderStatusRewarding
		// model.OrderStatusComplete
		return nil
	}

	notiMsg := msgqueue.AppNotiMsg{
		Category: msgqueue.NotiCategoryAffiliate,
		Title:    title,
		Body:     body,
		UserId:   userId,
		Data:     msgqueue.GetOrderUpdateNotiData(orderId),
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

	log.LG.Infof("Pushed Order Update Noti to queue %v\n", notiMsg)

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

func (u *RewardMaker) commitOrderUpdateMsg(message kafka.Message) error {
	err := u.orderUpdateQ.CommitMessages(context.Background(), message)
	if err != nil {
		log.LG.Errorf("Failed to commit order update message: %v", err)
		return err
	}
	return nil
}

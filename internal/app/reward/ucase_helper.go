package reward

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

func (u *RewardUsecase) CalculateWithdrawableReward(rewards []model.Reward, userId uint32) (*model.RewardWithdraw, []model.Reward, []model.OrderRewardHistory) {
	shippingRequestId := fmt.Sprintf("affiliate-%v:%v", userId, time.Now().UnixMilli())
	withdraw := model.RewardWithdraw{
		UserId:            uint(userId),
		ShippingRequestID: shippingRequestId,
		ShippingStatus:    model.ShippingStatusInit,
		Amount:            0,
		Fee:               AffRewardTxFee,
	}
	rewardsToWithdraw := []model.Reward{}
	orderRewardHistories := []model.OrderRewardHistory{}

	for idx := range rewards {
		orderReward, ended := rewards[idx].WithdrawableReward()
		if orderReward < MinWithdrawReward {
			continue
		}

		orderRewardHistories = append(orderRewardHistories, model.OrderRewardHistory{
			RewardID: rewards[idx].ID,
			Amount:   orderReward,
		})

		// Update total withdraw amount
		withdraw.Amount += orderReward

		// Update reward
		rewards[idx].RewardedAmount += orderReward
		if ended {
			rewards[idx].EndedAt = time.Now()
		}
		rewardsToWithdraw = append(rewardsToWithdraw, rewards[idx])
	}

	return &withdraw, rewardsToWithdraw, orderRewardHistories
}

package reward

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

func (u *rewardUCase) CalculateWithdrawalReward(rewards []model.Reward, userId uint32) (*model.RewardWithdraw, []model.Reward, []model.OrderRewardHistory) {
	shippingRequestId := fmt.Sprintf("affiliate-%v:%v", userId, time.Now().UnixMilli())
	withdraw := model.RewardWithdraw{
		UserId:            uint(userId),
		ShippingRequestID: shippingRequestId,
		ShippingStatus:    model.ShippingStatusInit,
		Amount:            0,
		Fee:               AffRewardTxFee,
	}
	var rewardsToWithdraw []model.Reward
	var orderRewardHistories []model.OrderRewardHistory

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

		if ended { // TODO: not update ended. Create status instead?
			rewards[idx].EndAt = time.Now()
		}
		rewardsToWithdraw = append(rewardsToWithdraw, rewards[idx])
	}

	return &withdraw, rewardsToWithdraw, orderRewardHistories
}

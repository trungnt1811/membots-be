package reward

import (
	"fmt"
	"math"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

func (u *RewardUsecase) calculateWithdrawableReward(rewards []model.Reward, userId uint32) (*model.RewardWithdraw, []model.Reward, []model.OrderRewardHistory) {
	shippingRequestId := fmt.Sprintf("affiliate-%v:%v", userId, time.Now().UnixMilli())
	withdraw := model.RewardWithdraw{
		UserId:            uint(userId),
		ShippingRequestID: shippingRequestId,
		Amount:            0,
		Fee:               AffRewardTxFee,
	}
	rewardsToWithdraw := []model.Reward{}
	orderRewardHistories := []model.OrderRewardHistory{}

	for idx := range rewards {
		orderReward, ended := rewards[idx].GetClaimableReward()
		if orderReward < MinWithdrawReward {
			continue
		}

		roundReward := math.Round(orderReward*100) / 100
		orderRewardHistories = append(orderRewardHistories, model.OrderRewardHistory{
			RewardID: rewards[idx].ID,
			Amount:   roundReward,
		})

		// Update total withdraw amount
		withdraw.Amount += roundReward

		// Update reward
		rewards[idx].RewardedAmount += roundReward
		if ended {
			rewards[idx].EndedAt = time.Now()
		}
		rewardsToWithdraw = append(rewardsToWithdraw, rewards[idx])
	}

	return &withdraw, rewardsToWithdraw, orderRewardHistories
}

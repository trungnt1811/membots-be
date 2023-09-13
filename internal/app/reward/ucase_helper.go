package reward

import (
	"fmt"
	"math"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

func (u *RewardUsecase) calculateClaimableReward(rewards []model.Reward, userId uint32) (*model.RewardClaim, []model.Reward, []model.OrderRewardHistory) {
	shippingRequestId := fmt.Sprintf("affiliate-%v:%v", userId, time.Now().UnixMilli())
	rewardClaim := model.RewardClaim{
		UserId:            uint(userId),
		ShippingRequestID: shippingRequestId,
		Amount:            0,
	}
	rewardToClaim := []model.Reward{}
	orderRewardHistories := []model.OrderRewardHistory{}

	for idx := range rewards {
		orderReward, ended := rewards[idx].GetClaimableReward()
		if orderReward < MinClaimReward {
			continue
		}

		roundReward := math.Round(orderReward*100) / 100
		orderRewardHistories = append(orderRewardHistories, model.OrderRewardHistory{
			RewardID: rewards[idx].ID,
			Amount:   roundReward,
		})

		// Update total claim amount
		rewardClaim.Amount += roundReward

		// Update reward
		rewards[idx].RewardedAmount += roundReward
		if ended {
			rewards[idx].EndedAt = time.Now()
		}
		rewardToClaim = append(rewardToClaim, rewards[idx])
	}

	return &rewardClaim, rewardToClaim, orderRewardHistories
}

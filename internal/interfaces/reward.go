package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward *model.Reward) error
	SaveRewardWithdraw(ctx context.Context, rewardClaim *model.RewardWithdraw, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory) error
	UpdateWithdrawShippingStatus(ctx context.Context, shippingReqId, txHash, status string) error
	GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardWithdraw, error)
	CountWithdrawal(ctx context.Context, userId uint32) (int64, error)
	GetTotalWithdrewAmount(ctx context.Context, userId uint32) (float64, error)
	GetInProgressRewards(ctx context.Context, userId uint32) ([]model.Reward, error)
	GetRewardsInDay(ctx context.Context) ([]model.Reward, error)
}

type RewardUCase interface {
	GetRewardSummary(ctx context.Context, userId uint32) (dto.RewardSummary, error)
	WithdrawReward(ctx context.Context, userId uint32, userWallet string) (dto.WithdrawRewardResponse, error)
	GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) (dto.RewardWithdrawResponse, error)
}

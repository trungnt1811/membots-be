package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward *model.Reward) error
	GetRewardById(ctx context.Context, userId uint, affOrderId uint) (model.Reward, error)
	GetRewardByOrderId(ctx context.Context, userId uint, affOrderId uint) (model.Reward, error)
	GetAllReward(ctx context.Context, userId uint, page, size int) ([]model.Reward, error)
	CountReward(ctx context.Context, userId uint) (int64, error)
	SaveRewardClaim(ctx context.Context, rewardClaim *model.RewardClaim, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory) error

	GetClaimHistory(ctx context.Context, userId uint, page, size int) ([]model.RewardClaim, error)
	CountClaimHistory(ctx context.Context, userId uint) (int64, error)
	GetInProgressRewards(ctx context.Context, userId uint) ([]model.Reward, error)
	GetRewardsInDay(ctx context.Context) ([]model.Reward, error)
}

type RewardUCase interface {
	GetRewardByOrderId(ctx context.Context, userId uint, affOrderId uint) (model.Reward, error)
	GetAllReward(ctx context.Context, userId uint, page, size int) (dto.RewardResponse, error)
	GetRewardSummary(ctx context.Context, userId uint) (dto.RewardSummary, error)
	ClaimReward(ctx context.Context, userId uint, userWallet string) (dto.ClaimRewardResponse, error)
	GetClaimHistory(ctx context.Context, userId uint, page, size int) (dto.RewardClaimResponse, error)
}

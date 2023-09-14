package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward *model.Reward) error
	GetRewardById(ctx context.Context, userId uint32, affOrderId uint) (model.Reward, error)
	GetRewardByOrderId(ctx context.Context, userId uint32, affOrderId uint) (model.Reward, error)
	GetAllReward(ctx context.Context, userId uint32, page, size int) ([]model.Reward, error)
	CountReward(ctx context.Context, userId uint32) (int64, error)
	SaveRewardClaim(ctx context.Context, rewardClaim *model.RewardClaim, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory) error

	GetClaimHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardClaim, error)
	CountClaimHistory(ctx context.Context, userId uint32) (int64, error)
	GetTotalClaimedAmount(ctx context.Context, userId uint32) (float64, error)
	GetInProgressRewards(ctx context.Context, userId uint32) ([]model.Reward, error)
	GetRewardsInDay(ctx context.Context) ([]model.Reward, error)
}

type RewardUCase interface {
	GetAllReward(ctx context.Context, userId uint32, page, size int) (dto.RewardResponse, error)
	GetRewardSummary(ctx context.Context, userId uint32) (dto.RewardSummary, error)
	ClaimReward(ctx context.Context, userId uint32, userWallet string) (dto.ClaimRewardResponse, error)
	GetClaimHistory(ctx context.Context, userId uint32, page, size int) (dto.RewardClaimResponse, error)
}

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

	GetRewardHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardHistoryFull, error)
	CountRewardHistory(ctx context.Context, userId uint32) (int64, error)
	GetInProgressRewards(ctx context.Context, userId uint32) ([]model.Reward, error)
	GetRewardedAmountByReward(ctx context.Context, rewardIds []uint) (map[uint]float64, error)
}

type RewardUCase interface {
	GetRewardByOrderId(ctx context.Context, userId uint32, affOrderId uint) (model.Reward, error)
	GetAllReward(ctx context.Context, userId uint32, page, size int) (dto.RewardResponse, error)
	GetRewardHistory(ctx context.Context, userId uint32, page, size int) (dto.RewardHistoryResponse, error)
	GetPendingRewards(ctx context.Context, userId uint32) ([]dto.RewardWithPendingDto, error)
}

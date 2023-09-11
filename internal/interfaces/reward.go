package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type RewardRepository interface {
	GetRewardByOrderId(ctx context.Context, affOrderId uint) (model.Reward, error)
	GetAllReward(ctx context.Context, userId uint, page, size int) ([]model.Reward, error)
	CountReward(ctx context.Context, userId uint) (int64, error)
	GetRewardHistory(ctx context.Context, userId uint, page, size int) ([]model.RewardHistoryFull, error)
	CountRewardHistory(ctx context.Context, userId uint) (int64, error)
}

type RewardUCase interface {
	GetRewardByOrderId(ctx context.Context, affOrderId uint) (model.Reward, error)
	GetAllReward(ctx context.Context, userId uint, page, size int) (dto.RewardResponse, error)
	GetRewardHistory(ctx context.Context, userId uint, page, size int) (dto.RewardHistoryResponse, error)
}

package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward *model.Reward) error
	GetRewardByAtOrderId(ctx context.Context, atOrderId string) (model.Reward, error)
	UpdateRewardByAtOrderId(atOrderId string, updates *model.Reward) error
	SaveRewardWithdraw(ctx context.Context, rewardClaim *model.RewardWithdraw, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory, completeRwOrders []string) error
	UpdateWithdrawShippingStatus(ctx context.Context, shippingReqId, txHash, status string) error
	GetWithdrawByShippingRequest(ctx context.Context, shippingReqId string) (model.RewardWithdraw, error)
	GetWithdrawById(ctx context.Context, userId uint32, withdrawId uint) (model.RewardWithdraw, error)
	GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardWithdraw, error)
	CountWithdrawal(ctx context.Context, userId uint32) (int64, error)
	GetTotalWithdrewAmount(ctx context.Context, userId uint32) (float64, error)
	GetInProgressRewards(ctx context.Context, userId uint32) ([]model.Reward, error)
	GetUsersHaveInProgressRewards(ctx context.Context) ([]uint32, error)
	GetInProgressRewardsOfMultipleUsers(ctx context.Context, userIds []uint32) ([]model.Reward, error)
}

type RewardUCase interface {
	GetRewardSummary(ctx context.Context, userId uint32) (dto.RewardSummary, error)
	WithdrawReward(ctx context.Context, userId uint32, userWallet string) (dto.RewardWithdrawDto, error)
	GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) (dto.RewardWithdrawResponse, error)
	GetWithdrawDetails(ctx context.Context, userId uint32, withdrawId uint) (dto.RewardWithdrawDto, error)
}

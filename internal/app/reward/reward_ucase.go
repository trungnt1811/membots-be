package reward

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type RewardUsecase struct {
	repo interfaces.RewardRepository
}

func NewRewardUsecase(repo interfaces.RewardRepository) *RewardUsecase {
	return &RewardUsecase{
		repo: repo,
	}
}

func (u *RewardUsecase) GetRewardByOrderId(ctx context.Context, affOrderId uint) (model.Reward, error) {
	return u.repo.GetRewardByOrderId(ctx, affOrderId)
}

func (u *RewardUsecase) GetAllReward(ctx context.Context, userId uint, page, size int) (dto.RewardResponse, error) {
	rewards, err := u.repo.GetAllReward(ctx, userId, page, size)
	if err != nil {
		return dto.RewardResponse{}, err
	}

	nextPage := page
	if len(rewards) > size {
		nextPage = page + 1
	}

	rewardDtos := make([]dto.RewardDto, len(rewards))
	for idx, item := range rewards {
		rewardDtos[idx] = item.ToRewardDto()
	}

	totalReward, err := u.repo.CountReward(ctx, userId)
	if err != nil {
		return dto.RewardResponse{}, err
	}

	return dto.RewardResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     rewardDtos,
		Total:    totalReward,
	}, nil
}

func (u *RewardUsecase) GetRewardHistory(ctx context.Context, userId uint, page, size int) (dto.RewardHistoryResponse, error) {
	rewards, err := u.repo.GetRewardHistory(ctx, userId, page, size)
	if err != nil {
		return dto.RewardHistoryResponse{}, err
	}

	nextPage := page
	if len(rewards) > size {
		nextPage = page + 1
	}

	rewardDtos := make([]dto.RewardHistoryDto, len(rewards))
	for idx, item := range rewards {
		rewardDtos[idx] = item.ToRewardHistoryDto()
	}

	totalRewardHistory, err := u.repo.CountRewardHistory(ctx, userId)
	if err != nil {
		return dto.RewardHistoryResponse{}, err
	}

	return dto.RewardHistoryResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     rewardDtos,
		Total:    totalRewardHistory,
	}, nil
}

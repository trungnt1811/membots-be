package reward

import (
	"context"
	"fmt"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/shipping"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const (
	MinClaimReward         = 0.01    // asa
	RewardLockTime         = 60 * 24 // hours
	AffCommissionFee       = 0       // percentage preserve for stella from affiliate reward
	AffRewardTxFee         = 0.1     // fix amount of tx fee
	StellaSellerId         = 119
	StellaAffRewardProgram = "0x4744fd21dbA951890c0247b6585930E74AC6A3d5"
)

type RewardUsecase struct {
	repo      interfaces.RewardRepository
	orderRepo interfaces.OrderRepository
	rwService *shipping.ShippingClient
}

func NewRewardUsecase(repo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository,
	rwService *shipping.ShippingClient,
) *RewardUsecase {
	return &RewardUsecase{
		repo:      repo,
		orderRepo: orderRepo,
		rwService: rwService,
	}
}

func (u *RewardUsecase) ClaimReward(ctx context.Context, userId uint32, userWallet string) (dto.ClaimRewardResponse, error) {
	rewards, err := u.repo.GetInProgressRewards(ctx, userId)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	// Calculating Reward
	rewardClaim, rewardToClaim, orderRewardHistories := u.calculateClaimableReward(rewards, userId)
	if rewardClaim.Amount-AffRewardTxFee < MinClaimReward {
		return dto.ClaimRewardResponse{
			Execute: false,
			Amount:  rewardClaim.Amount,
		}, nil
	}

	// Call service send reward
	sendReq := shipping.ReqSendPayload{
		SellerId:       StellaSellerId,
		ProgramAddress: StellaAffRewardProgram,
		RequestId:      rewardClaim.ShippingRequestID,
		Items: []shipping.ReqSendItem{
			{
				WalletAddress: userWallet,
				Amount:        fmt.Sprint(rewardClaim.Amount - AffRewardTxFee),
			},
		},
	}
	_, err = u.rwService.SendReward(&sendReq)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	// Update Db
	err = u.repo.SaveRewardClaim(ctx, rewardClaim, rewardToClaim, orderRewardHistories)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	return dto.ClaimRewardResponse{
		Execute: true,
		Amount:  rewardClaim.Amount,
	}, nil
}

func (u *RewardUsecase) GetRewardSummary(ctx context.Context, userId uint32) (dto.RewardSummary, error) {
	totalClaimedAmount, err := u.repo.GetTotalClaimedAmount(ctx, userId)
	if err != nil {
		return dto.RewardSummary{}, err
	}

	inProgressRewards, err := u.repo.GetInProgressRewards(ctx, userId)
	if err != nil {
		return dto.RewardSummary{}, err
	}

	var totalOrderReward float64 = 0
	for _, item := range inProgressRewards {
		totalOrderReward += item.Amount - item.RewardedAmount
	}

	rewardsInDay, err := u.repo.GetRewardsInDay(ctx)
	if err != nil {
		return dto.RewardSummary{}, err
	}
	var totalOrderRewardInDay float64 = 0
	for _, item := range rewardsInDay {
		totalOrderRewardInDay += item.Amount * model.FirstPartRewardPercent
	}

	rewardClaim, _, _ := u.calculateClaimableReward(inProgressRewards, userId)

	return dto.RewardSummary{
		TotalClaimedAmount:  totalClaimedAmount,
		ClaimableReward:     rewardClaim.Amount,
		RewardInDay:         totalOrderRewardInDay,
		PendingRewardOrder:  len(inProgressRewards),
		PendingRewardAmount: totalOrderReward - rewardClaim.Amount,
	}, nil
}

func (u *RewardUsecase) GetAllReward(ctx context.Context, userId uint32, page, size int) (dto.RewardResponse, error) {
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

func (u *RewardUsecase) GetClaimHistory(ctx context.Context, userId uint32, page, size int) (dto.RewardClaimResponse, error) {
	rewards, err := u.repo.GetClaimHistory(ctx, userId, page, size)
	if err != nil {
		return dto.RewardClaimResponse{}, err
	}

	nextPage := page
	if len(rewards) > size {
		nextPage = page + 1
	}

	rewardDtos := make([]dto.RewardClaimDto, len(rewards))
	for idx, item := range rewards {
		rewardDtos[idx] = item.ToRewardClaimDto()
	}

	totalRewardHistory, err := u.repo.CountClaimHistory(ctx, userId)
	if err != nil {
		return dto.RewardClaimResponse{}, err
	}

	return dto.RewardClaimResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     rewardDtos,
		Total:    totalRewardHistory,
	}, nil
}

func (u *RewardUsecase) GetClaimDetails(ctx context.Context, userId uint32, claimId uint) (dto.RewardClaimDetailsDto, error) {
	return u.repo.GetClaimDetails(ctx, userId, claimId)
}

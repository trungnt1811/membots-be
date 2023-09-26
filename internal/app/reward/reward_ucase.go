package reward

import (
	"context"
	"fmt"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/shipping"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
)

const (
	MinWithdrawReward = 0.01    // asa
	RewardLockTime    = 60 * 24 // hours
	AffRewardTxFee    = 0.1     // fix amount of tx fee
)

type RewardConfig struct {
	RewardProgram string // Reward Program that send affiliate reward
	SellerId      uint   // Owner of Reward Program
}

type rewardUCase struct {
	repo         interfaces.RewardRepository
	orderRepo    interfaces.OrderRepository
	rwService    *shipping.ShippingClient
	rewardConfig RewardConfig
}

func NewRewardUCase(repo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository,
	rwService *shipping.ShippingClient,
	rewardConfig RewardConfig,
) interfaces.RewardUCase {
	return &rewardUCase{
		repo:         repo,
		orderRepo:    orderRepo,
		rwService:    rwService,
		rewardConfig: rewardConfig,
	}
}

func (u *rewardUCase) WithdrawReward(ctx context.Context, userId uint32, userWallet string) (dto.WithdrawRewardResponse, error) {
	rewards, err := u.repo.GetInProgressRewards(ctx, userId)
	if err != nil {
		return dto.WithdrawRewardResponse{}, err
	}

	// Calculating Reward
	rewardClaim, rewardToClaim, orderRewardHistories, completeRwOrders := u.CalculateWithdrawalReward(rewards, userId)
	if rewardClaim.Amount-AffRewardTxFee < MinWithdrawReward {
		return dto.WithdrawRewardResponse{
			Execute: false,
			Amount:  rewardClaim.Amount,
		}, nil
	}

	// Call service send reward
	sendReq := shipping.ReqSendPayload{
		SellerId:       u.rewardConfig.SellerId,
		ProgramAddress: u.rewardConfig.RewardProgram,
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
		return dto.WithdrawRewardResponse{}, err
	}

	// Save Db
	rewardClaim.ShippingStatus = model.ShippingStatusSending
	err = u.repo.SaveRewardWithdraw(ctx, rewardClaim, rewardToClaim, orderRewardHistories, completeRwOrders)
	if err != nil {
		return dto.WithdrawRewardResponse{}, err
	}

	return dto.WithdrawRewardResponse{
		Execute: true,
		Amount:  rewardClaim.Amount,
	}, nil
}

func (u *rewardUCase) GetRewardSummary(ctx context.Context, userId uint32) (dto.RewardSummary, error) {
	totalWithdrewAmount, err := u.repo.GetTotalWithdrewAmount(ctx, userId)
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
		totalOrderRewardInDay += item.Amount * item.ImmediateRelease
	}
	totalOrderRewardInDay = util.RoundFloat(totalOrderRewardInDay, 2)

	withdrawable, _, _, _ := u.CalculateWithdrawalReward(inProgressRewards, userId)
	pendingRewardAmount := util.RoundFloat(totalOrderReward-withdrawable.Amount, 2)

	return dto.RewardSummary{
		TotalWithdrewAmount: totalWithdrewAmount,
		WithdrawableReward:  withdrawable.Amount,
		RewardInDay:         totalOrderRewardInDay,
		PendingRewardOrder:  len(inProgressRewards),
		PendingRewardAmount: pendingRewardAmount,
	}, nil
}

func (u *rewardUCase) GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) (dto.RewardWithdrawResponse, error) {
	rewards, err := u.repo.GetWithdrawHistory(ctx, userId, page, size)
	if err != nil {
		return dto.RewardWithdrawResponse{}, err
	}

	nextPage := page
	if len(rewards) > size {
		nextPage = page + 1
	}

	rewardDtos := make([]dto.RewardWithdrawDto, len(rewards))
	for idx, item := range rewards {
		rewardDtos[idx] = item.ToRewardWithdrawDto()
	}

	totalRewardHistory, err := u.repo.CountWithdrawal(ctx, userId)
	if err != nil {
		return dto.RewardWithdrawResponse{}, err
	}

	return dto.RewardWithdrawResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     rewardDtos,
		Total:    totalRewardHistory,
	}, nil
}

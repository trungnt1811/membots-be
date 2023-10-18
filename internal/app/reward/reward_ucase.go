package reward

import (
	"context"
	"fmt"
	"time"

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

func (u *rewardUCase) WithdrawReward(ctx context.Context, userId uint32, userWallet string) (dto.RewardWithdrawDto, error) {
	rewards, err := u.repo.GetInProgressRewards(ctx, userId)
	if err != nil {
		return dto.RewardWithdrawDto{}, err
	}

	// Calculating Reward
	rewardClaim, rewardToClaim, orderRewardHistories, completeRwOrders := u.CalculateWithdrawalReward(rewards, userId)
	if rewardClaim.Amount-AffRewardTxFee < MinWithdrawReward {
		return dto.RewardWithdrawDto{}, fmt.Errorf("not enough reward to withdraw, minimum %v ASA", AffRewardTxFee+MinWithdrawReward)
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
		return dto.RewardWithdrawDto{}, err
	}

	// Save Db
	rewardClaim.ShippingStatus = model.ShippingStatusSending
	err = u.repo.SaveRewardWithdraw(ctx, rewardClaim, rewardToClaim, orderRewardHistories, completeRwOrders)
	if err != nil {
		return dto.RewardWithdrawDto{}, err
	}

	return rewardClaim.ToRewardWithdrawDto(), nil
}

func (u *rewardUCase) GetRewardSummary(ctx context.Context, userId uint32) (dto.RewardSummary, error) {
	totalWithdrewAmount, err := u.repo.GetTotalWithdrewAmount(ctx, userId)
	if err != nil {
		return dto.RewardSummary{}, err
	}

	pastTimeLimit := time.Now().Add(-6 * 30 * 24 * time.Hour)
	waitForApproveOrder, err := u.orderRepo.CountOrders(ctx, pastTimeLimit, userId, dto.OrderStatusWaitForConfirming)
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

	var totalOrderRewardInDay float64 = 0
	now := time.Now()
	for _, r := range inProgressRewards {
		if now.Before(r.EndAt) && now.After(r.StartAt.Add(model.OneDay)) {
			totalOrderRewardInDay += r.OneDayReward()
		}
	}

	withdrawable, _, _, _ := u.CalculateWithdrawalReward(inProgressRewards, userId)
	pendingRewardAmount := util.RoundFloat(totalOrderReward-withdrawable.Amount, 2)

	return dto.RewardSummary{
		TotalWithdrewAmount: util.RoundFloat(totalWithdrewAmount, 2),
		WithdrawableReward:  util.RoundFloat(withdrawable.Amount, 2),
		RewardInDay:         util.RoundFloat(totalOrderRewardInDay, 2),
		PendingRewardOrder:  len(inProgressRewards),
		PendingRewardAmount: pendingRewardAmount,
		WaitForApproveOrder: waitForApproveOrder,
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

func (u *rewardUCase) GetWithdrawDetails(ctx context.Context, userId uint32, withdrawId uint) (dto.RewardWithdrawDto, error) {
	withdraw, err := u.repo.GetWithdrawById(ctx, userId, withdrawId)
	if err != nil {
		return dto.RewardWithdrawDto{}, err
	}

	return withdraw.ToRewardWithdrawDto(), nil
}

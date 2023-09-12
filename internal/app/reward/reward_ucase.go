package reward

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/infra/shipping"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/samber/lo"
)

const (
	MinClaimReward         = 0.01    // asa
	RewardLockTime         = 60 * 24 // hours
	AffRewardFee           = 0       // percentage preserve for stella from affiliate reward
	StellaSellerId         = 119
	StellaAffRewardProgram = "0x4744fd21dbA951890c0247b6585930E74AC6A3d5"
)

type RewardUsecase struct {
	repo      interfaces.RewardRepository
	orderRepo interfaces.OrderRepository
	rwService *shipping.ShippingClient
	approveQ  *msgqueue.QueueReader
}

func NewRewardUsecase(repo interfaces.RewardRepository,
	orderRepo interfaces.OrderRepository,
	rwService *shipping.ShippingClient,
	approveQ *msgqueue.QueueReader) *RewardUsecase {
	return &RewardUsecase{
		repo:      repo,
		orderRepo: orderRepo,
		rwService: rwService,
		approveQ:  approveQ,
	}
}

func (u *RewardUsecase) ListenOrderApproved() {
	for {
		ctx := context.Background()

		newAtOrderId, err := u.getNewApprovedAtOrderId(ctx)
		if err != nil {
			log.Println("Failed to comsume new AccessTrade order id", err)
			continue
		}

		order, err := u.orderRepo.FindOrderByAccessTradeId(newAtOrderId)
		if err != nil {
			log.Println("Failed to get affiliate order", err)
			continue
		}

		now := time.Now()
		newReward := model.Reward{
			UserId:         order.UserId,
			AtOrderID:      newAtOrderId,
			Amount:         float64(order.PubCommission) * AffRewardFee / 100,
			RewardedAmount: 0,
			Fee:            AffRewardFee,
			EndedAt:        now.Add(RewardLockTime * time.Hour),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		err = u.repo.CreateReward(ctx, &newReward)
		if err != nil {
			log.Println("Failed to get affiliate order", err)
			continue
		}

		// TODO: send notification ?
	}
}

func (u *RewardUsecase) ClaimReward(ctx context.Context, userId uint, userWallet string, rewardId uint) (dto.ClaimRewardResponse, error) {
	reward, err := u.repo.GetRewardById(ctx, userId, rewardId)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}
	claimableReward := reward.GetClaimableReward()
	if claimableReward < MinClaimReward {
		return dto.ClaimRewardResponse{
			Execute: false,
			Amount:  claimableReward,
		}, nil
	}
	// TODO: claim all reward, not by reward id
	sendReq := shipping.ReqSendPayload{
		SellerId:       StellaSellerId,
		ProgramAddress: StellaAffRewardProgram,
		RequestId:      fmt.Sprintf("affiliate-%v-%v:%v", reward.ID, time.Now().UnixMilli(), reward.AtOrderID),
		Items: []shipping.ReqSendItem{
			{
				WalletAddress: userWallet,
				Amount:        fmt.Sprint(claimableReward),
			},
		},
	}
	_, err = u.rwService.SendReward(&sendReq)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	return dto.ClaimRewardResponse{
		Execute: true,
		Amount:  claimableReward,
	}, nil
}

func (u *RewardUsecase) GetRewardByOrderId(ctx context.Context, userId uint, affOrderId uint) (model.Reward, error) {
	return u.repo.GetRewardByOrderId(ctx, userId, affOrderId)
}

func (u *RewardUsecase) GetPendingRewards(ctx context.Context, userId uint) ([]dto.RewardWithPendingDto, error) {
	inProgressReward, err := u.repo.GetInProgressRewards(ctx, userId)
	if err != nil {
		return []dto.RewardWithPendingDto{}, err
	}
	rewardIds := lo.Map(inProgressReward, func(item model.Reward, _ int) uint {
		return item.ID
	})
	rewardedReward, err := u.repo.GetRewardedAmountByReward(ctx, rewardIds)
	if err != nil {
		return []dto.RewardWithPendingDto{}, err
	}

	rewardWithPending := make([]dto.RewardWithPendingDto, len(inProgressReward))
	for idx, item := range inProgressReward {
		rewardWithPending[idx] = dto.RewardWithPendingDto{
			ID:            item.ID,
			PendingAmount: item.Amount - rewardedReward[item.ID],
		}
	}

	return rewardWithPending, nil
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

package reward

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/infra/shipping"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const (
	MinClaimReward         = 0.001   // asa
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

func (u *RewardUsecase) ClaimReward(ctx context.Context, userId uint, userWallet string) (dto.ClaimRewardResponse, error) {
	rewards, err := u.repo.GetInProgressRewards(ctx, userId)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	// Calculating Reward
	shippingRequestId := fmt.Sprintf("affiliate-%v:%v", userId, time.Now().UnixMilli())
	rewardClaim := model.RewardClaim{
		UserId:            userId,
		ShippingRequestID: shippingRequestId,
		Amount:            0,
	}
	rewardToClaim := []model.Reward{}
	orderRewardHistories := []model.OrderRewardHistory{}

	for idx := range rewards {
		orderReward, ended := rewards[idx].GetClaimableReward()
		if orderReward < MinClaimReward {
			continue
		}

		roundReward := math.Round(orderReward*100) / 100
		orderRewardHistories = append(orderRewardHistories, model.OrderRewardHistory{
			RewardID: rewards[idx].ID,
			Amount:   roundReward,
		})

		// Update total claim amount
		rewardClaim.Amount += roundReward

		// Update reward
		rewards[idx].RewardedAmount += roundReward
		if ended {
			rewards[idx].EndedAt = time.Now()
		}
		rewardToClaim = append(rewardToClaim, rewards[idx])
	}

	if len(rewardToClaim) == 0 {
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
				Amount:        fmt.Sprint(rewardClaim.Amount),
			},
		},
	}
	_, err = u.rwService.SendReward(&sendReq)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	// Update Db
	err = u.repo.SaveRewardClaim(ctx, &rewardClaim, rewardToClaim, orderRewardHistories)
	if err != nil {
		return dto.ClaimRewardResponse{}, err
	}

	return dto.ClaimRewardResponse{
		Execute: true,
		Amount:  rewardClaim.Amount,
	}, nil
}

func (u *RewardUsecase) GetRewardByOrderId(ctx context.Context, userId uint, affOrderId uint) (model.Reward, error) {
	return u.repo.GetRewardByOrderId(ctx, userId, affOrderId)
}

func (u *RewardUsecase) GetRewardSummary(ctx context.Context, userId uint) (dto.RewardSummary, error) {
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
		ClaimableReward:     rewardClaim.Amount,
		RewardInDay:         totalOrderRewardInDay,
		PendingRewardOrder:  len(inProgressRewards),
		PendingRewardAmount: totalOrderReward - rewardClaim.Amount,
	}, nil
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

func (u *RewardUsecase) GetClaimHistory(ctx context.Context, userId uint, page, size int) (dto.RewardClaimResponse, error) {
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

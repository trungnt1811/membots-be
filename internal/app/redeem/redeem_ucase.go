package redeem

import (
	"github.com/astraprotocol/affiliate-system/internal/app/redeem/types"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

const (
	AVAILABLE_SECONDS = 300 // 5min
)

type RedeemUsecase struct {
	repo interfaces.RedeemRepository
}

func NewRedeemUsecase(repo interfaces.RedeemRepository) *RedeemUsecase {
	return &RedeemUsecase{
		repo: repo,
	}
}

func (rd *RedeemUsecase) RedeemCashback(req types.RedeemRequestPayload) (*types.RedeemRewardResponse, error) {
	return nil, nil
}

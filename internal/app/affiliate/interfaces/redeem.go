package interfaces

import "github.com/astraprotocol/affiliate-system/internal/app/affiliate/redeem/types"

type RedeemRepository interface {
	TopUpCashBack(userId int, amount float32) error
}

type RedeemUsecase interface {
	RedeemCashback(types.RedeemRequestPayload) (*types.RedeemRewardResponse, error)
}

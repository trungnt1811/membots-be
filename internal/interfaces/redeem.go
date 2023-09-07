package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/app/redeem/types"
)

type RedeemRepository interface {
	TopUpCashBack(userId int, amount float32) error
}

type RedeemUCase interface {
	RedeemCashback(types.RedeemRequestPayload) (*types.RedeemRewardResponse, error)
}

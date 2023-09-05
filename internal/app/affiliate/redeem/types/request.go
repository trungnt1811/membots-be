package types

import (
	"errors"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util/validators"
	"github.com/ethereum/go-ethereum/common"
)

type RedeemRequestPayload struct {
	// a unique text to redeem reward
	RedeemCode string `form:"redeemCode" json:"redeemCode" binding:"required"`
	// valid wallet address to send reward to
	WalletAddress string `form:"walletAddress" json:"walletAddress" binding:"required"`
}

func (payload *RedeemRequestPayload) Valid() error {
	err := validators.ValidateAddress(payload.WalletAddress)
	if err != nil {
		return errors.New("wallet address malform")
	}

	return nil
}

type CheckRedeemRequest struct {
	// a unique text to redeem reward
	RedeemCode string `form:"redeemCode" json:"redeemCode" binding:"required"`

	// TODO: check if bot spam
}

type ClaimSuccessPayload struct {
	RedeemCode  string `form:"redeemCode" json:"redeemCode" binding:"required"`
	ClaimAt     int64  `form:"claimAt" json:"claimAt" binding:"required"`
	ClaimTxHash string `form:"claimTxHash" json:"claimTxHash"`
}

func (payload *ClaimSuccessPayload) Valid() error {
	if payload.ClaimAt <= 0 {
		return errors.New("claim at must be seconds > 0")
	}
	txHash := common.HexToHash(payload.ClaimTxHash[2:])
	if txHash.Hex() == common.HexToHash("0").Hex() {
		return errors.New("wrong tx hash format")
	}

	return nil
}

type ClaimSuccessResponse struct {
}

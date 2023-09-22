package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

const (
	ShippingStatusInit    = "initial"
	ShippingStatusSending = "sending"
	ShippingStatusSuccess = "success"
	ShippingStatusFailed  = "failed"
)

func (m *RewardWithdraw) TableName() string {
	return "aff_reward_withdraw"
}

type RewardWithdraw struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	UserId            uint      `json:"user_id"`
	ShippingRequestID string    `json:"shipping_request_id"`
	TxHash            string    `json:"tx_hash"`
	ShippingStatus    string    `json:"shipping_status"`
	Amount            float64   `json:"amount"` // amount before tx fee subtraction
	Fee               float64   `json:"fee"`    // withdraw fee, prevent dos
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (r *RewardWithdraw) ToRewardWithdrawDto() dto.RewardWithdrawDto {
	return dto.RewardWithdrawDto{
		ID:                r.ID,
		UserId:            r.UserId,
		ShippingRequestID: r.ShippingRequestID,
		Amount:            r.Amount,
		Fee:               r.Fee,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

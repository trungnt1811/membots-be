package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

func (m *RewardClaim) TableName() string {
	return "aff_reward_claim"
}

type RewardClaim struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	UserId            uint      `json:"user_id"`
	ShippingRequestID string    `json:"shipping_request_id"`
	Amount            float64   `json:"amount"` // amount before tx fee subtraction
	Fee               float64   `json:"fee"`    // tx fee
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (r *RewardClaim) ToRewardClaimDto() dto.RewardClaimDto {
	return dto.RewardClaimDto{
		ID:                r.ID,
		UserId:            r.UserId,
		ShippingRequestID: r.ShippingRequestID,
		Amount:            r.Amount,
		Fee:               r.Fee,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

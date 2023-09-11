package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

func (m *Reward) TableName() string {
	return "aff_reward"
}

type Reward struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserId     uint      `json:"user_id"`
	AffOrderID uint      `json:"aff_order_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *Reward) ToRewardDto() dto.RewardDto {
	return dto.RewardDto{
		ID:         r.ID,
		UserId:     r.UserId,
		AffOrderID: r.AffOrderID,
		Amount:     r.Amount,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}

}

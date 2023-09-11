package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

const (
	RewardStatusInProgress = "in_progress"
	RewardStatusDone       = "done"
)

func (m *Reward) TableName() string {
	return "aff_reward"
}

type RewardedByReward struct {
	RewardID       uint    `json:"reward_id"`
	RewardedAmount float64 `json:"rewarded_amount"`
}

type Reward struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserId     uint      `json:"user_id"`
	AffOrderID uint      `json:"aff_order_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	EndedAt    time.Time `json:"ended_at"`
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
		EndedAt:    r.CreatedAt,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

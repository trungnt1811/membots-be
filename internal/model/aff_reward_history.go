package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

func (m *RewardHistory) TableName() string {
	return "aff_reward_history"
}

type RewardHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RewardID  uint      `json:"reward_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RewardHistoryFull struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	RewardID   uint      `json:"reward_id"`
	UserId     uint      `json:"user_id"`
	AffOrderID uint      `json:"aff_order_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *RewardHistoryFull) ToRewardHistoryDto() dto.RewardHistoryDto {
	return dto.RewardHistoryDto{
		ID:         r.ID,
		RewardID:   r.RewardID,
		UserId:     r.UserId,
		AffOrderID: r.AffOrderID,
		Amount:     r.Amount,
		Type:       r.Type,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

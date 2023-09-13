package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

const (
	RewardTypeFirstTime = "first_time"
	RewardTypeClaim     = "claim"
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
	ID        uint      `gorm:"primarykey" json:"id"`
	RewardID  uint      `json:"reward_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uint      `json:"user_id"`
	AtOrderID string    `json:"accesstrade_order_id"`
}

func (r *RewardHistoryFull) ToRewardHistoryDto() dto.RewardHistoryDto {
	return dto.RewardHistoryDto{
		ID:        r.ID,
		RewardID:  r.RewardID,
		Amount:    r.Amount,
		Type:      r.Type,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		UserId:    r.UserId,
		AtOrderID: r.AtOrderID,
	}
}

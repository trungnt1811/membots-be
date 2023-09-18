package model

import (
	"time"
)

func (m *OrderRewardHistory) TableName() string {
	return "aff_reward_history"
}

type OrderRewardHistory struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	RewardID         uint      `json:"reward_id"`
	RewardWithdrawID uint      `json:"reward_withdraw_id"`
	Amount           float64   `json:"amount"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

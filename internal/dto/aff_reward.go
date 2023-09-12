package dto

import (
	"time"
)

// REWARD

type RewardDto struct {
	ID             uint      `json:"id"`
	UserId         uint      `json:"user_id"`
	AtOrderID      string    `json:"accesstrade_order_id"`
	Amount         float64   `json:"amount"`
	RewardedAmount float64   `json:"rewarded_amount"`
	EndedAt        time.Time `json:"ended_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RewardWithPendingDto struct {
	ID            uint    `json:"id"`
	PendingAmount float64 `json:"pending_amount"`
}

type ClaimRewardResponse struct {
	Execute bool    `json:"execute"`
	Amount  float64 `json:"amount"`
}

type RewardResponse struct {
	NextPage int         `json:"next_page"`
	Page     int         `json:"page"`
	Size     int         `json:"size"`
	Data     []RewardDto `json:"data"`
	Total    int64       `json:"total"`
}

// REWARD HISTORY

type RewardHistoryDto struct {
	ID        uint      `json:"id"`
	RewardID  uint      `json:"reward_id"`
	UserId    uint      `json:"user_id"`
	AtOrderID string    `json:"accesstrade_order_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RewardHistoryResponse struct {
	NextPage int                `json:"next_page"`
	Page     int                `json:"page"`
	Size     int                `json:"size"`
	Data     []RewardHistoryDto `json:"data"`
	Total    int64              `json:"total"`
}

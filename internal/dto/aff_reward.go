package dto

import (
	"time"
)

// REWARD

type RewardDto struct {
	ID         uint      `json:"id"`
	UserId     uint      `json:"user_id"`
	AffOrderID uint      `json:"aff_order_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	EndedAt    time.Time `json:"ended_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RewardWithPendingDto struct {
	ID            uint    `json:"id"`
	PendingAmount float64 `json:"pending_amount"`
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
	ID               uint      `json:"id"`
	RewardID         uint      `json:"reward_id"`
	UserId           uint      `json:"user_id"`
	AffOrderID       uint      `json:"aff_order_id"`
	Amount           float64   `json:"amount"`
	CumulativeAmount float64   `json:"cumulative_amount"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RewardHistoryResponse struct {
	NextPage int                `json:"next_page"`
	Page     int                `json:"page"`
	Size     int                `json:"size"`
	Data     []RewardHistoryDto `json:"data"`
	Total    int64              `json:"total"`
}

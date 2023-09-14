package dto

import (
	"time"
)

// REWARD
type RewardSummary struct {
	ClaimableReward     float64 `json:"claimable_reward"`
	RewardInDay         float64 `json:"reward_in_day"`
	PendingRewardOrder  int     `json:"pending_reward_order"`
	PendingRewardAmount float64 `json:"pending_reward_amount"`
}

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

type RewardClaimDto struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	UserId            uint      `json:"user_id"`
	ShippingRequestID string    `json:"shipping_request_id"`
	Amount            float64   `json:"amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RewardClaimResponse struct {
	NextPage int              `json:"next_page"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
	Data     []RewardClaimDto `json:"data"`
	Total    int64            `json:"total"`
}

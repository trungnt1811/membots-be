package dto

import (
	"time"
)

// REWARD
type RewardSummary struct {
	// total amount of reward already withdrawn
	TotalWithdrewAmount float64 `json:"total_withdrew_amount"`
	// amount of reward earn in day = 50% of total order approved today
	RewardInDay float64 `json:"reward_in_day"`
	// number of orders that user not receive all reward
	PendingRewardOrder int `json:"pending_reward_order"`
	// amount of reward can be withdrawn at the moment
	WithdrawableReward float64 `json:"withdrawable_reward"`
	// amount of reward still locked
	PendingRewardAmount float64 `json:"pending_reward_amount"`
	// number of orders still waiting for approve
	WaitForApproveOrder int64 `json:"wait_for_approve_order"`
}

type RewardDto struct {
	ID             uint      `json:"id"`
	UserId         uint      `json:"user_id"`
	AtOrderID      string    `json:"accesstrade_order_id"`
	Amount         float64   `json:"amount"`
	RewardedAmount float64   `json:"rewarded_amount"`
	CommissionFee  float64   `json:"commission_fee"` // commission fee (in percentage)
	EndAt          time.Time `json:"end_at"`
	StartAt        time.Time `json:"start_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RewardResponse struct {
	NextPage int         `json:"next_page"`
	Page     int         `json:"page"`
	Size     int         `json:"size"`
	Data     []RewardDto `json:"data"`
	Total    int64       `json:"total"`
}

// REWARD HISTORY

type RewardWithdrawDto struct {
	ID             uint      `json:"id"`
	UserId         uint      `json:"user_id"`
	TxHash         string    `json:"tx_hash"`
	ShippingStatus string    `json:"shipping_status"`
	Amount         float64   `json:"amount"`
	Fee            float64   `json:"fee"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RewardWithdrawResponse struct {
	NextPage int                 `json:"next_page"`
	Page     int                 `json:"page"`
	Size     int                 `json:"size"`
	Data     []RewardWithdrawDto `json:"data"`
	Total    int64               `json:"total"`
}

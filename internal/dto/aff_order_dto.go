package dto

import "time"

const (
	OrderStatusWaitForConfirming = "wait_for_confirming"
	OrderStatusRewarding         = "rewarding"
	OrderStatusComplete          = "complete"
	OrderStatusCancelled         = "cancelled"
	OrderStatusRejected          = "rejected"
)

type OrderDetailsDto struct {
	ID                 uint      `json:"id"`
	UserId             uint      `json:"user_id"`
	OrderStatus        string    `json:"order_status"`
	ATProductLink      string    `json:"at_product_link"`
	Billing            float32   `json:"billing"`
	CategoryName       string    `json:"category_name"`
	Merchant           string    `json:"merchant"`
	AccessTradeOrderId string    `json:"accesstrade_order_id"`
	PubCommission      float32   `json:"pub_commission"`
	SalesTime          time.Time `json:"sales_time"`     // dat hang thanh cong
	ConfirmedTime      time.Time `json:"confirmed_time"` // don hang duoc ghi nhan
	RewardAmount       float64   `json:"amount"`         // amount of reward after fee subtractions
	RewardedAmount     float64   `json:"rewarded_amount"`
	CommissionFee      float64   `json:"commission_fee"` // commission fee (in percentage)
	RewardEndAt        time.Time `json:"reward_end_at"`
	RewardStartAt      time.Time `json:"reward_start_at"` // xac nhan hoan tat nhan 50%
}

type OrderHistoryResponse struct {
	NextPage int               `json:"next_page"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	Data     []OrderDetailsDto `json:"data"`
	Total    int64             `json:"total"`
}

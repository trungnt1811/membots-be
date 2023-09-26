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
	ID                             uint      `json:"id"`
	UserId                         uint      `json:"user_id"`
	OrderStatus                    string    `json:"order_status"`
	ATProductLink                  string    `json:"at_product_link"`
	Billing                        float32   `json:"billing"`
	CategoryName                   string    `json:"category_name"`
	Merchant                       string    `json:"merchant"`
	ImageUrl                       string    `json:"image_url"` // brand logo
	AccessTradeOrderId             string    `json:"accesstrade_order_id"`
	PubCommission                  float32   `json:"pub_commission"`
	CommissionFee                  float64   `json:"commission_fee"`                            // commission fee (in percentage)
	BuyTime                        time.Time `json:"sales_time,omitempty"`                      // dat hang thanh cong
	ConfirmedTime                  time.Time `json:"confirmed_time,omitempty"`                  // don hang duoc ghi nhan
	RejectedTime                   time.Time `json:"rejected_time,omitempty"`                   // don hang bi huy
	RewardFirstPartReleasedTime    time.Time `json:"reward_first_part_released_time,omitempty"` // xac nhan hoan tat nhan 50%
	RewardFirstPartReleasedAmount  float64   `json:"reward_first_part_released_amount,omitempty"`
	RewardSecondPartUnlockedAmount float64   `json:"reward_second_part_unlocked_amount,omitempty"`
	RewardRemainingAmount          float64   `json:"reward_remaining_amount,omitempty"`
	RewardAmount                   float64   `json:"reward_amount,omitempty"`
	RewardEndAt                    time.Time `json:"reward_end_at,omitempty"`
}

type OrderHistoryResponse struct {
	NextPage int               `json:"next_page"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	Data     []OrderDetailsDto `json:"data"`
	Total    int64             `json:"total"`
}

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
	Billing                        float32   `json:"billing"`
	CategoryName                   string    `json:"category_name"`
	Merchant                       string    `json:"merchant"`
	Brand                          BrandDto  `json:"brand"`
	AccessTradeOrderId             string    `json:"accesstrade_order_id"`
	PubCommission                  float32   `json:"pub_commission"`
	CommissionFee                  float64   `json:"commission_fee"`                  // commission fee (in percentage)
	BuyTime                        time.Time `json:"sales_time"`                      // dat hang thanh cong
	ConfirmedTime                  time.Time `json:"confirmed_time"`                  // don hang duoc ghi nhan
	RejectedTime                   time.Time `json:"rejected_time"`                   // don hang bi tu choi
	CancelledTime                  time.Time `json:"cancelled_time"`                  // don hang bi huy
	RewardFirstPartReleasedTime    time.Time `json:"reward_first_part_released_time"` // xac nhan hoan tat nhan 50%
	RewardFirstPartReleasedAmount  float64   `json:"reward_first_part_released_amount"`
	RewardSecondPartUnlockedAmount float64   `json:"reward_second_part_unlocked_amount"`
	RewardRemainingAmount          float64   `json:"reward_remaining_amount"`
	RewardAmount                   float64   `json:"reward_amount"`
	RewardEndAt                    time.Time `json:"reward_end_at"`
}

type OrderHistoryResponse struct {
	NextPage int               `json:"next_page"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	Data     []OrderDetailsDto `json:"data"`
	Total    int64             `json:"total"`
}

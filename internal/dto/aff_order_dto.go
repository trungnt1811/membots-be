package dto

import "time"

type OrderDetailsDto struct {
	UserId             uint      `json:"user_id"`
	OrderStatus        string    `json:"order_status"`
	ATProductLink      string    `json:"at_product_link"`
	Billing            float32   `json:"billing"`
	CategoryName       string    `json:"category_name"`
	Merchant           string    `json:"merchant"`
	AccessTradeOrderId string    `json:"accesstrade_order_id"`
	PubCommission      float32   `json:"pub_commission"`
	SalesTime          time.Time `json:"sales_time"`
	ConfirmedTime      time.Time `json:"confirmed_time"`
	ApprovedTime       time.Time `json:"approved_time"`
	RejectedTime       time.Time `json:"rejected_time"`
	RewardAmount       float64   `json:"amount"` // amount of reward after fee subtractions
	RewardedAmount     float64   `json:"rewarded_amount"`
	CommissionFee      float64   `json:"commission_fee"` // commission fee (in percentage)
	EndedAt            time.Time `json:"ended_at"`
	CreatedAt          time.Time `json:"created_at"`
}

type OrderHistoryResponse struct {
	NextPage int               `json:"next_page"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	Data     []OrderDetailsDto `json:"data"`
	Total    int64             `json:"total"`
}

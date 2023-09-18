package dto

import "time"

type OrderDetailsDto struct {
	UserId             uint                 `json:"user_id"`
	OrderStatus        string               `json:"order_status"`
	ATProductLink      string               `json:"at_product_link"`
	Billing            float32              `json:"billing"`
	CategoryName       string               `json:"category_name"`
	ConfirmedTime      time.Time            `json:"confirmed_time"`
	Merchant           string               `json:"merchant"`
	AccessTradeOrderId string               `json:"accesstrade_order_id"`
	PubCommission      float32              `json:"pub_commission"`
	SalesTime          time.Time            `json:"sales_time"`
	Timeline           map[string]time.Time `json:"timeline"` // status changing history
	Reward             *RewardDto           `json:"reward"`
}

type OrderHistoryResponse struct {
	NextPage int               `json:"next_page"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	Data     []OrderDetailsDto `json:"data"`
	Total    int64             `json:"total"`
}

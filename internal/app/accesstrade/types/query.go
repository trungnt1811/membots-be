package types

import "time"

const (
	ORDER_STATUS_PENDING  = 0
	ORDER_STATUS_APPROVED = 1
	ORDER_STATUS_REJECTED = 2
)

type ATOrderQuery struct {
	Since       time.Time `url:"since"`
	Until       time.Time `url:"until"`
	Status      uint8     `url:"status,omitempty"`
	Merchant    string    `url:"merchant,omitempty"`
	UTMSource   string    `url:"utm_source,omitempty"`
	UTMCampaign string    `url:"utm_campaign,omitempty"`
	UTMMedium   string    `url:"utm_medium,omitempty"`
	UTMContent  string    `url:"utm_content,omitempty"`
}

type ATTransactionQuery struct {
	Since         time.Time `url:"since"`
	Until         time.Time `url:"until"`
	Status        uint8     `url:"status,omitempty"`
	IsConfirmed   uint8     `url:"is_confirmed,omitempty"`
	TransactionId string    `url:"transaction_id,omitempty"`
	UpdateTimeEnd time.Time `url:"update_time_end,omitempty"`
	Merchant      string    `url:"merchant,omitempty"`
	UTMSource     string    `url:"utm_source,omitempty"`
	UTMCampaign   string    `url:"utm_campaign,omitempty"`
	UTMMedium     string    `url:"utm_medium,omitempty"`
	UTMContent    string    `url:"utm_content,omitempty"`
}

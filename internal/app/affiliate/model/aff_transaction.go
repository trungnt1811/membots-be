package model

import "time"

type AffTransaction struct {
	ID                       uint           `gorm:"primarykey" json:"id"`
	CampaignId               uint           `json:"campaign_id"`
	UserId                   uint           `json:"user_id"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	AccessTradeId            string         `json:"accesstrade_id"`
	AccessTradeTransactionId string         `json:"accesstrade_transaction_id"`
	AccessTradeConversionId  string         `json:"accesstrade_conversion_id"`
	Merchant                 string         `json:"merchant"`
	Status                   uint8          `json:"status"`
	ClickTime                time.Time      `json:"click_time"`
	TransactionTime          time.Time      `json:"transaction_time"`
	TransactionValue         float32        `json:"transaction_value"`
	UpdateTime               time.Time      `json:"update_time"`
	ConfirmedTime            time.Time      `json:"confirmed_time"`
	IsConfirmed              uint8          `json:"is_confirmed"`
	Commission               float32        `json:"commission"`
	ProductId                string         `json:"product_id"`
	ProductName              string         `json:"product_name"`
	ProductPrice             string         `json:"product_price"`
	ProductQuantity          int            `json:"product_quantity"`
	ProductImage             string         `json:"product_image"`
	ProductCategory          string         `json:"product_category"`
	Extra                    map[string]any `json:"extra"`
	CategoryName             string         `json:"category_name"`
	ConversionPlatform       string         `json:"conversion_platform"`
	ClickURL                 string         `json:"click_url"`
	UTMTerm                  string         `json:"utm_term"`
	UTMSource                string         `json:"utm_source"`
	UTMCampaign              string         `json:"utm_campaign"`
	UTMMedium                string         `json:"utm_medium"`
	UTMContent               string         `json:"utm_content"`
	ReasonRejected           string         `json:"reason_rejected"`
	CustomerType             string         `json:"customer_type"`
}

func (m *AffTransaction) TableName() string {
	return "aff_transaction"
}

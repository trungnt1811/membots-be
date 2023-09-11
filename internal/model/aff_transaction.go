package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
)

type AffTransaction struct {
	ID                      uint           `gorm:"primarykey" json:"id"`
	CampaignId              uint           `json:"campaign_id"`
	UserId                  uint           `json:"user_id"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	AccessTradeId           string         `json:"accesstrade_id"`
	AccessTradeOrderId      string         `json:"accesstrade_order_id"` // NOTE: AccessTrade Order ID
	AccessTradeConversionId int64          `json:"accesstrade_conversion_id"`
	Merchant                string         `json:"merchant"`
	Status                  uint8          `json:"status"`
	ClickTime               time.Time      `json:"click_time"`
	TransactionTime         time.Time      `json:"transaction_time"`
	TransactionValue        float32        `json:"transaction_value"`
	UpdateTime              time.Time      `json:"update_time"`
	ConfirmedTime           time.Time      `json:"confirmed_time"`
	IsConfirmed             uint8          `json:"is_confirmed"`
	Commission              float32        `json:"commission"`
	ProductId               string         `json:"product_id"`
	ProductName             string         `json:"product_name"`
	ProductPrice            float32        `json:"product_price"`
	ProductQuantity         int            `json:"product_quantity"`
	ProductImage            string         `json:"product_image"`
	ProductCategory         string         `json:"product_category"`
	Extra                   map[string]any `json:"extra"`
	CategoryName            string         `json:"category_name"`
	ConversionPlatform      string         `json:"conversion_platform"`
	ClickURL                string         `json:"click_url"`
	UTMTerm                 string         `json:"utm_term"`
	UTMSource               string         `json:"utm_source"`
	UTMCampaign             string         `json:"utm_campaign"`
	UTMMedium               string         `json:"utm_medium"`
	UTMContent              string         `json:"utm_content"`
	ReasonRejected          string         `json:"reason_rejected"`
	CustomerType            string         `json:"customer_type"`
}

func (m *AffTransaction) TableName() string {
	return "aff_transaction"
}

func NewAffTransactionFromAT(order *AffOrder, atTx *types.ATTransaction) *AffTransaction {
	return &AffTransaction{
		UserId:                  order.UserId,
		Merchant:                atTx.Merchant,
		Status:                  atTx.Status,
		UpdateTime:              atTx.ClickTime.Time,
		ClickURL:                atTx.ClickUrl,
		ConversionPlatform:      atTx.ConversionPlatform,
		UTMTerm:                 atTx.UTMTerm,
		UTMSource:               atTx.UTMSource,
		UTMCampaign:             atTx.UTMCampaign,
		UTMContent:              atTx.UTMContent,
		UTMMedium:               atTx.UTMMedium,
		ProductCategory:         atTx.ProductCategory,
		TransactionTime:         atTx.TransactionTime.Time,
		ProductImage:            atTx.ProductImage,
		TransactionValue:        atTx.TransactionValue,
		Extra:                   atTx.Extra,
		ReasonRejected:          atTx.ReasonRejected,
		CategoryName:            atTx.CategoryName,
		ProductId:               atTx.ProductId,
		IsConfirmed:             atTx.IsConfirmed,
		ConfirmedTime:           atTx.ConfirmedTime.Time,
		ProductPrice:            atTx.ProductPrice,
		AccessTradeId:           atTx.Id,
		Commission:              atTx.Commission,
		CustomerType:            atTx.CustomerType,
		AccessTradeConversionId: atTx.ConversionId,
		ProductQuantity:         atTx.ProductQuantity,
		ClickTime:               atTx.ClickTime.Time,
		ProductName:             atTx.ProductName,
		AccessTradeOrderId:      atTx.TransactionId,
	}
}

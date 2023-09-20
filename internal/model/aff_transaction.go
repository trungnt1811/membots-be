package model

import (
	"encoding/json"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"gorm.io/datatypes"
)

type AffTransaction struct {
	ID                      uint           `gorm:"primarykey" json:"id"`
	UserId                  uint           `json:"user_id"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	AccessTradeId           string         `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	AccessTradeOrderId      string         `json:"accesstrade_order_id" gorm:"column:accesstrade_order_id"` // NOTE: AccessTrade Order ID
	AccessTradeConversionId int64          `json:"accesstrade_conversion_id" gorm:"column:accesstrade_conversion_id"`
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
	Extra                   datatypes.JSON `json:"extra"`
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

func (m *AffTransaction) ToDto() dto.AffTransaction {
	var extra map[string]any
	err := json.Unmarshal([]byte(m.Extra.String()), &extra)
	if err != nil {
		log.LG.Errorf("scan _extra error: %v", err)
	}
	return dto.AffTransaction{
		ID:                      m.ID,
		UserId:                  m.UserId,
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
		AccessTradeId:           m.AccessTradeId,
		AccessTradeOrderId:      m.AccessTradeOrderId,
		AccessTradeConversionId: m.AccessTradeConversionId,
		Merchant:                m.Merchant,
		Status:                  m.Status,
		ClickTime:               m.ClickTime,
		TransactionTime:         m.TransactionTime,
		TransactionValue:        m.TransactionValue,
		UpdateTime:              m.UpdateTime,
		ConfirmedTime:           m.ConfirmedTime,
		IsConfirmed:             m.IsConfirmed,
		Commission:              m.Commission,
		ProductId:               m.ProductId,
		ProductName:             m.ProductName,
		ProductPrice:            m.ProductPrice,
		ProductQuantity:         m.ProductQuantity,
		ProductImage:            m.ProductImage,
		ProductCategory:         m.ProductCategory,
		Extra:                   extra,
		CategoryName:            m.CategoryName,
		ConversionPlatform:      m.ConversionPlatform,
		ClickURL:                m.ClickURL,
		UTMTerm:                 m.UTMTerm,
		UTMSource:               m.UTMSource,
		UTMCampaign:             m.UTMCampaign,
		UTMMedium:               m.UTMMedium,
		UTMContent:              m.UTMContent,
		ReasonRejected:          m.ReasonRejected,
		CustomerType:            m.CustomerType,
	}
}

func NewAffTransactionFromAT(order *AffOrder, atTx *types.ATTransaction) *AffTransaction {
	extraBytes, err := json.Marshal(atTx.Extra)
	if err != nil {
		log.LG.Errorf("cannot unmarshal extra: %v", err)
		extraBytes = []byte("{}")
	}
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
		Extra:                   datatypes.JSON(extraBytes),
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

package dto

import (
	"time"
)

type TimeRange struct {
	Since time.Time `json:"since" form:"since"` // Time in RFC3339 format
	Until time.Time `json:"until" form:"until"` // Time in RFC3339 format
}

type OrderListQuery struct {
	Page        int       `form:"page" json:"page"`
	PerPage     int       `form:"per_page" json:"per_page"`
	UserId      uint      `form:"user_id" json:"user_id"`
	OrderStatus string    `form:"order_status" json:"order_status"`
	Since       time.Time `form:"since" json:"since"`
	Until       time.Time `form:"until" json:"until"`
}

type OrderListResponse struct {
	Data       []AffOrder `json:"data"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}

type PostBackListQuery struct {
	Page    int       `form:"page" json:"page"`
	PerPage int       `form:"per_page" json:"per_page"`
	OrderId string    `form:"order_id" json:"order_id"`
	IsError bool      `form:"is_error" json:"is_error"`
	Since   time.Time `form:"since" json:"since"`
	Until   time.Time `form:"until" json:"until"`
}

type PostBackListResponse struct {
	Data       []AffPostBack `json:"data"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

type AffTransaction struct {
	ID                      uint      `gorm:"primarykey" json:"id"`
	UserId                  uint      `json:"user_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	AccessTradeId           string    `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	AccessTradeOrderId      string    `json:"accesstrade_order_id" gorm:"column:accesstrade_order_id"` // NOTE: AccessTrade Order ID
	AccessTradeConversionId int64     `json:"accesstrade_conversion_id" gorm:"column:accesstrade_conversion_id"`
	Merchant                string    `json:"merchant"`
	Status                  uint8     `json:"status"`
	ClickTime               time.Time `json:"click_time"`
	TransactionTime         time.Time `json:"transaction_time"`
	TransactionValue        float32   `json:"transaction_value"`
	UpdateTime              time.Time `json:"update_time"`
	ConfirmedTime           time.Time `json:"confirmed_time"`
	IsConfirmed             uint8     `json:"is_confirmed"`
	Commission              float32   `json:"commission"`
	ProductId               string    `json:"product_id"`
	ProductName             string    `json:"product_name"`
	ProductPrice            float32   `json:"product_price"`
	ProductQuantity         int       `json:"product_quantity"`
	ProductImage            string    `json:"product_image"`
	ProductCategory         string    `json:"product_category"`
	Extra                   any       `json:"extra"`
	CategoryName            string    `json:"category_name"`
	ConversionPlatform      string    `json:"conversion_platform"`
	ClickURL                string    `json:"click_url"`
	UTMTerm                 string    `json:"utm_term"`
	UTMSource               string    `json:"utm_source"`
	UTMCampaign             string    `json:"utm_campaign"`
	UTMMedium               string    `json:"utm_medium"`
	UTMContent              string    `json:"utm_content"`
	ReasonRejected          string    `json:"reason_rejected"`
	CustomerType            string    `json:"customer_type"`
}

type AffOrder struct {
	ID                 uint             `gorm:"primarykey" json:"id"`
	AffLink            string           `json:"aff_link"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	UserId             uint             `json:"user_id"`
	OrderStatus        string           `json:"order_status"`
	ATProductLink      string           `json:"at_product_link"`
	Billing            float32          `json:"billing"`
	Browser            string           `json:"browser"`
	CategoryName       string           `json:"category_name"`
	ClientPlatform     string           `json:"client_platform"`
	ClickTime          time.Time        `json:"click_time"`
	ConfirmedTime      time.Time        `json:"confirmed_time"`
	ConversionPlatform string           `json:"conversion_platform"`
	CustomerType       string           `json:"customer_type"`
	IsConfirmed        uint8            `json:"is_confirmed"`
	LandingPage        string           `json:"landing_page"`
	Merchant           string           `json:"merchant"`
	AccessTradeOrderId string           `json:"accesstrade_order_id" gorm:"column:accesstrade_order_id"`
	OrderPending       uint8            `json:"order_pending"`
	OrderReject        uint8            `json:"order_reject"`
	OrderApproved      uint8            `json:"order_approved"`
	ProductCategory    string           `json:"product_category"`
	ProductsCount      int              `json:"products_count"`
	PubCommission      float32          `json:"pub_commission"`
	SalesTime          time.Time        `json:"sales_time"`
	UpdateTime         time.Time        `json:"update_time"`
	Website            string           `json:"website"`
	WebsiteURL         string           `json:"website_url"`
	UTMTerm            string           `json:"utm_term"`
	UTMSource          string           `json:"utm_source"`
	UTMCampaign        string           `json:"utm_campaign"`
	UTMMedium          string           `json:"utm_medium"`
	UTMContent         string           `json:"utm_content"`
	Transactions       []AffTransaction `json:"transactions"`
}

type AffPostBack struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	OrderId      string         `json:"order_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Data         map[string]any `json:"data"`
	ErrorMessage string         `json:"error_message"`
}

type SyncOrderRewardPayload struct {
	AccessTradeOrderId string `json:"accesstrade_order_id"`
}

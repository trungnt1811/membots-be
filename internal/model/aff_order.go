package model

import "time"

type AffOrder struct {
	ID                 uint      `gorm:"primarykey" json:"id"`
	CampaignId         uint      `json:"campaign_id"`
	AffLink            string    `json:"aff_link"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	UserId             uint      `json:"user_id"`
	OrderStatus        string    `json:"order_status"`
	ATProductLink      string    `json:"at_product_link"`
	Billing            float32   `json:"billing"`
	Browser            string    `json:"browser"`
	CategoryName       string    `json:"category_name"`
	ClientPlatform     string    `json:"client_platform"`
	ClickTime          time.Time `json:"click_time"`
	ConfirmedTime      time.Time `json:"confirmed_time"`
	ConversionPlatform string    `json:"conversion_platform"`
	CustomerType       string    `json:"customer_type"`
	IsConfirmed        uint8     `json:"is_confirmed"`
	LandingPage        string    `json:"landing_page"`
	Merchant           string    `json:"merchant"`
	AccessTradeOrderId string    `json:"accesstrade_order_id"`
	OrderPending       uint8     `json:"order_pending"`
	OrderReject        uint8     `json:"order_reject"`
	OrderApproved      uint8     `json:"order_approved"`
	ProductCategory    string    `json:"product_category"`
	ProductCount       string    `json:"products_count"`
	PubCommission      float32   `json:"pub_commission"`
	SalesTime          time.Time `json:"sales_time"`
	UpdateTime         time.Time `json:"update_time"`
	Website            string    `json:"website"`
	WebsiteURL         string    `json:"website_url"`
	UTMTerm            string    `json:"utm_term"`
	UTMSource          string    `json:"utm_source"`
	UTMCampaign        string    `json:"utm_campaign"`
	UTMMedium          string    `json:"utm_medium"`
	UTMContent         string    `json:"utm_content"`
}

func (m *AffOrder) TableName() string {
	return "aff_order"
}

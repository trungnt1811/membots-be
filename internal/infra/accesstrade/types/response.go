package types

import (
	"fmt"
	"strings"
	"time"
)

const CUSTOM_DAY_LAYOUT = "2022-08-08"

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	// Split and parse
	t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s+07:00", value))

	if err != nil {
		return err
	}
	*c = CustomTime{t} // set result using the pointer
	return nil
}

func (c CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.Local().Format(time.RFC3339)[0:19] + `"`), nil
}

type ATCampaignDescription struct {
	ActionPoint           string `json:"action_point"`
	CommissionPolicy      string `json:"commission_policy"`
	CookiePolicy          string `json:"cookie_policy"`
	Introduction          string `json:"introduction"`
	OtherNotice           string `json:"other_notice"`
	RejectedReason        string `json:"rejected_reason"`
	TrafficBuildingPolicy string `json:"traffic_building_policy"`
}

type ATCampaign struct {
	Id             string                `json:"id"`
	Logo           string                `json:"logo"`
	MaxCom         string                `json:"max_com"`
	Merchant       string                `json:"merchant"`
	Name           string                `json:"name"`
	Scope          string                `json:"scope"`
	Approval       string                `json:"approval"`
	Status         int                   `json:"status"`
	Type           int                   `json:"type"`
	Url            string                `json:"url"`
	Category       string                `json:"category"`
	SubCategory    string                `json:"sub_category"`
	CookieDuration int                   `json:"cookie_duration"`
	CookiePolicy   string                `json:"cookie_policy"`
	Description    ATCampaignDescription `json:"description"`
	StartTime      CustomTime            `json:"start_time"`
	EndTime        CustomTime            `json:"end_time"`
}

type ATCampaignListResp struct {
	Total     int          `json:"total"`
	Data      []ATCampaign `json:"data"`
	Page      int          `json:"page"`
	TotalPage float32      `json:"total_page"`
}

type ErrorLink struct {
	CampaignId   string `json:"campaign_id"`
	CampaignName string `json:"campaign_name"`
	Message      string `json:"message"`
	UrlOrigin    string `json:"url_origin"`
}

type ATLink struct {
	AffLink   string `json:"aff_link"`
	FirstLink string `json:"first_link"`
	ShortLink string `json:"short_link"`
	UrlOrigin string `json:"url_origin"`
}

type AllLink struct {
	ErrorLink   []ErrorLink `json:"error_link"`
	SuccessLink []ATLink    `json:"success_link"`
	SuspendUrl  []any       `json:"suspend_url"`
}

type ATLinkResp struct {
	Data       AllLink `json:"data"`
	Success    bool    `json:"success"`
	Message    string  `json:"message"`
	StatusCode int     `json:"status_code"`
}

type ATTransaction struct {
	Id                 string         `json:"id"`
	ConversionId       int64          `json:"conversion_id"`
	Merchant           string         `json:"merchant"`
	Status             uint8          `json:"status"`
	ClickTime          CustomTime     `json:"click_time"`
	TransactionId      string         `json:"transaction_id"`
	TransactionTime    CustomTime     `json:"transaction_time"`
	TransactionValue   float32        `json:"transaction_value"`
	UpdateTime         CustomTime     `json:"update_time"`
	ConfirmedTime      CustomTime     `json:"confirmed_time"`
	IsConfirmed        uint8          `json:"is_confirmed"`
	Commission         float32        `json:"commission"`
	ProductId          string         `json:"product_id"`
	ProductName        string         `json:"product_name"`
	ProductPrice       float32        `json:"product_price"`
	ProductQuantity    int            `json:"product_quantity"`
	ProductImage       string         `json:"product_image"`
	ProductCategory    string         `json:"product_category"`
	Extra              map[string]any `json:"_extra"`
	CategoryName       string         `json:"category_name"`
	ConversionPlatform string         `json:"conversion_platform"`
	ClickUrl           string         `json:"click_url"`
	UTMTerm            string         `json:"utm_term"`
	UTMSource          string         `json:"utm_source"`
	UTMMedium          string         `json:"utm_medium"`
	UTMCampaign        string         `json:"utm_campaign"`
	UTMContent         string         `json:"utm_content"`
	ReasonRejected     string         `json:"reason_rejected"`
	CustomerType       string         `json:"customer_type"`
}

type ATTransactionResp struct {
	Total     int             `json:"total"`
	Data      []ATTransaction `json:"data"`
	Page      int             `json:"page"`
	TotalPage int             `json:"total_page"`
}

type ATOrder struct {
	ATProductLink      string     `json:"at_product_link"`
	Billing            float32    `json:"billing"`
	Browser            string     `json:"browser"`
	CategoryName       string     `json:"category_name"`
	ClientPlatform     string     `json:"client_platform"`
	ClickTime          CustomTime `json:"click_time"`
	ConfirmedTime      CustomTime `json:"confirmed_time"`
	ConversionPlatform string     `json:"conversion_platform"`
	CustomerType       string     `json:"customer_type"`
	IsConfirmed        uint8      `json:"is_confirmed"`
	LandingPage        string     `json:"landing_page"`
	Merchant           string     `json:"merchant"`
	OrderId            string     `json:"order_id"`
	OrderApproved      uint       `json:"order_approved"`
	OrderPending       uint       `json:"order_pending"`
	OrderReject        uint       `json:"order_reject"`
	ProductCategory    string     `json:"product_category"`
	ProductsCount      int        `json:"products_count"`
	PubCommission      float32    `json:"pub_commission"`
	SalesTime          CustomTime `json:"sales_time"`
	UpdateTime         CustomTime `json:"update_time"`
	UTMSource          string     `json:"utm_source"`
	UTMMedium          string     `json:"utm_medium"`
	UTMCampaign        string     `json:"utm_campaign"`
	UTMContent         string     `json:"utm_content"`
	Website            string     `json:"website"`
	WebsiteUrl         string     `json:"website_url"`
}

type ATOrderListResp struct {
	Total     int       `json:"total"`
	Data      []ATOrder `json:"data"`
	Page      int       `json:"page"`
	TotalPage int       `json:"total_page"`
}

type ATMerchant struct {
	DisplayName []string `json:"display_name"`
	Id          string   `json:"id"`
	LoginName   string   `json:"login_name"`
	Logo        string   `json:"logo"`
	TotalOffer  int      `json:"total_offer"`
}

type ATMerchantResp struct {
	Data []ATMerchant `json:"data"`
}

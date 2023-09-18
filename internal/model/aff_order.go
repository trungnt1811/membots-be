package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
)

type AffOrder struct {
	ID                 uint      `gorm:"primarykey" json:"id"`
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
	AccessTradeOrderId string    `json:"accesstrade_order_id" gorm:"column:accesstrade_order_id"`
	OrderPending       uint8     `json:"order_pending"`
	OrderReject        uint8     `json:"order_reject"`
	OrderApproved      uint8     `json:"order_approved"`
	ProductCategory    string    `json:"product_category"`
	ProductsCount      int       `json:"products_count"`
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

func (order *AffOrder) TableName() string {
	return "aff_order"
}

func NewOrderFromATOrder(userId uint, atOrder *types.ATOrder) *AffOrder {
	orderStatus := "initial"
	if atOrder.OrderPending != 0 {
		orderStatus = "pending"
	} else if atOrder.OrderApproved != 0 {
		orderStatus = "approved"
	} else if atOrder.OrderReject != 0 {
		orderStatus = "rejected"
	}

	return &AffOrder{
		UserId:             userId,
		OrderStatus:        orderStatus,
		ATProductLink:      atOrder.ATProductLink,
		Billing:            atOrder.Billing,
		Browser:            atOrder.Browser,
		CategoryName:       atOrder.CategoryName,
		ClientPlatform:     atOrder.ClientPlatform,
		ClickTime:          atOrder.ClickTime.Time,
		ConfirmedTime:      atOrder.ConfirmedTime.Time,
		ConversionPlatform: atOrder.ConversionPlatform,
		CustomerType:       atOrder.CustomerType,
		IsConfirmed:        atOrder.IsConfirmed,
		LandingPage:        atOrder.LandingPage,
		Merchant:           atOrder.Merchant,
		AccessTradeOrderId: atOrder.OrderId,
		OrderApproved:      atOrder.OrderApproved,
		OrderPending:       atOrder.OrderPending,
		OrderReject:        atOrder.OrderReject,
		ProductCategory:    atOrder.ProductCategory,
		ProductsCount:      atOrder.ProductsCount,
		PubCommission:      atOrder.PubCommission,
		SalesTime:          atOrder.SalesTime.Time,
		UpdateTime:         atOrder.UpdateTime.Time,
		UTMSource:          atOrder.UTMSource,
		UTMCampaign:        atOrder.UTMCampaign,
		UTMMedium:          atOrder.UTMMedium,
		UTMContent:         atOrder.UTMContent,
		Website:            atOrder.Website,
		WebsiteURL:         atOrder.WebsiteUrl,
	}
}

type OrderDetails struct {
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
	Reward             *Reward              `json:"reward"`
}

func (o *OrderDetails) ToOrderDetailsDto() dto.OrderDetailsDto {
	reward := o.Reward.ToRewardDto()
	return dto.OrderDetailsDto{
		UserId:             o.UserId,
		OrderStatus:        o.OrderStatus,
		ATProductLink:      o.ATProductLink,
		Billing:            o.Billing,
		CategoryName:       o.CategoryName,
		ConfirmedTime:      o.ConfirmedTime,
		Merchant:           o.Merchant,
		AccessTradeOrderId: o.AccessTradeOrderId,
		PubCommission:      o.PubCommission,
		SalesTime:          o.SalesTime,
		Timeline:           o.Timeline,
		Reward:             &reward,
	}
}

// func (order *AffOrder) UpdateFromATOrder(req *dto.ATPostBackRequest) {
// 	order.ProductCategory = req.ProductCategory
// 	order.ProductCount = fmt.Sprint(req.Quantity)
// 	switch req.Status {
// 	case dto.REQ_STATUS_NEW:
// 		order.OrderStatus = "pending"
// 		order.OrderPending = 1
// 		order.OrderReject = 0
// 		order.OrderApproved = 0
// 		break
// 	case dto.REQ_STATUS_APPROVED:
// 		order.OrderStatus = "approved"
// 		order.OrderPending = 0
// 		order.OrderReject = 0
// 		order.OrderApproved = 1
// 		break
// 	case dto.REQ_STATUS_REJECTED:
// 		order.OrderStatus = "rejected"
// 		order.OrderPending = 0
// 		order.OrderReject = 1
// 		order.OrderApproved = 0
// 		break
// 	default:
// 		break
// 	}
// 	if parsed, err := util.ParsePostBackTime(req.SalesTime); err == nil {
// 		order.SalesTime = parsed
// 	}
// 	if parsed, err := util.ParsePostBackTime(req.ClickTime); err == nil {
// 		order.ClickTime = parsed
// 	}
// 	order.Browser = req.Browser
// 	order.UTMSource = req.UTMSource
// 	order.UTMCampaign = req.UTMCampaign
// 	order.UTMContent = req.UTMContent
// 	order.UTMMedium = req.UTMMedium
// 	order.IsConfirmed = req.IsConfirmed
// 	order.CustomerType = req.CustomerType
// 	order.UpdatedAt = time.Now()
// }

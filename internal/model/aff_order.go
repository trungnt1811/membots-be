package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/util"
)

const (
	OrderStatusInitial   = "initial"
	OrderStatusPending   = "pending"
	OrderStatusApproved  = "approved"
	OrderStatusRejected  = "rejected"
	OrderStatusCancelled = "cancelled"
	OrderStatusRewarding = "rewarding"
	OrderStatusComplete  = "complete"
)

type AffOrder struct {
	ID                 uint       `gorm:"primarykey" json:"id"`
	AffLink            string     `json:"aff_link"`
	CampaignId         uint       `json:"campaign_id"`
	BrandId            uint       `json:"brand_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UserId             uint       `json:"user_id"`
	OrderStatus        string     `json:"order_status"`
	ATProductLink      string     `json:"at_product_link"`
	Billing            float32    `json:"billing"`
	Browser            string     `json:"browser"`
	CategoryName       string     `json:"category_name"`
	ClientPlatform     string     `json:"client_platform"`
	ClickTime          time.Time  `json:"click_time"`
	ConfirmedTime      time.Time  `json:"confirmed_time"`
	ConversionPlatform string     `json:"conversion_platform"`
	CustomerType       string     `json:"customer_type"`
	IsConfirmed        uint8      `json:"is_confirmed"`
	LandingPage        string     `json:"landing_page"`
	Merchant           string     `json:"merchant"`
	AccessTradeOrderId string     `json:"accesstrade_order_id" gorm:"column:accesstrade_order_id"`
	OrderPending       *uint8     `json:"order_pending"`
	OrderReject        *uint8     `json:"order_reject"`
	OrderApproved      *uint8     `json:"order_approved"`
	ProductCategory    string     `json:"product_category"`
	ProductsCount      int        `json:"products_count"`
	PubCommission      float32    `json:"pub_commission"`
	SalesTime          time.Time  `json:"sales_time"`
	UpdateTime         time.Time  `json:"update_time"`
	CancelledTime      *time.Time `json:"cancelled_time"`
	Website            string     `json:"website"`
	WebsiteURL         string     `json:"website_url"`
	UTMTerm            string     `json:"utm_term"`
	UTMSource          string     `json:"utm_source"`
	UTMCampaign        string     `json:"utm_campaign"`
	UTMMedium          string     `json:"utm_medium"`
	UTMContent         string     `json:"utm_content"`
}

func (order *AffOrder) TableName() string {
	return "aff_order"
}

func GetAffOrderStatusFromAtOrder(atOrder *types.ATOrder) string {
	orderStatus := OrderStatusInitial
	if atOrder.OrderPending != 0 {
		orderStatus = OrderStatusPending
	}
	if atOrder.OrderApproved != 0 && atOrder.OrderPending == 0 {
		orderStatus = OrderStatusApproved
	}
	if atOrder.OrderReject != 0 {
		orderStatus = OrderStatusRejected
	}
	return orderStatus
}

func NewOrderFromATOrder(userId uint, campaignId uint, brandId uint, atOrder *types.ATOrder) *AffOrder {
	return &AffOrder{
		UserId:             userId,
		CampaignId:         campaignId,
		BrandId:            brandId,
		OrderStatus:        GetAffOrderStatusFromAtOrder(atOrder),
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
		OrderApproved:      &atOrder.OrderApproved,
		OrderPending:       &atOrder.OrderPending,
		OrderReject:        &atOrder.OrderReject,
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
		CancelledTime:      nil,
	}
}

func (old *AffOrder) CheckStatusChanged(newUpdate *AffOrder) bool {
	if old.OrderStatus == OrderStatusInitial {
		// already notified to user, not send anymore
		if newUpdate.OrderStatus == OrderStatusPending {
			return false
		}
	}

	if old.OrderStatus == OrderStatusRewarding {
		if newUpdate.OrderStatus == OrderStatusApproved {
			return false
		}
	}

	if old.OrderStatus == OrderStatusComplete {
		return false
	}

	return true
}

func (o *AffOrder) ToDto() dto.AffOrder {
	return dto.AffOrder{
		ID:                 o.ID,
		AffLink:            o.AffLink,
		CreatedAt:          o.CreatedAt,
		UpdatedAt:          o.UpdatedAt,
		UserId:             o.UserId,
		OrderStatus:        o.OrderStatus,
		ATProductLink:      o.ATProductLink,
		Billing:            o.Billing,
		Browser:            o.Browser,
		CategoryName:       o.CategoryName,
		ClientPlatform:     o.ClientPlatform,
		ClickTime:          o.ClickTime,
		ConfirmedTime:      o.ConfirmedTime,
		ConversionPlatform: o.ConversionPlatform,
		CustomerType:       o.CustomerType,
		IsConfirmed:        o.IsConfirmed,
		LandingPage:        o.LandingPage,
		Merchant:           o.Merchant,
		AccessTradeOrderId: o.AccessTradeOrderId,
		OrderPending:       *o.OrderPending,
		OrderReject:        *o.OrderReject,
		OrderApproved:      *o.OrderApproved,
		ProductCategory:    o.ProductCategory,
		ProductsCount:      o.ProductsCount,
		PubCommission:      o.PubCommission,
		SalesTime:          o.SalesTime,
		UpdateTime:         o.UpdateTime,
		Website:            o.Website,
		WebsiteURL:         o.WebsiteURL,
		UTMTerm:            o.UTMTerm,
		UTMSource:          o.UTMSource,
		UTMCampaign:        o.UTMCampaign,
		UTMMedium:          o.UTMMedium,
		UTMContent:         o.UTMContent,
	}
}

type OrderDetails struct {
	ID                 uint      `json:"id"`
	UserId             uint      `json:"user_id"`
	OrderStatus        string    `json:"order_status"`
	Billing            float32   `json:"billing"`
	CategoryName       string    `json:"category_name"`
	Merchant           string    `json:"merchant"`
	AccessTradeOrderId string    `json:"accesstrade_order_id"`
	PubCommission      float32   `json:"pub_commission"`
	UpdateTime         time.Time `json:"update_time"`
	CancelledTime      time.Time `json:"cancelled_time"`
	SalesTime          time.Time `json:"sales_time"`
	ConfirmedTime      time.Time `json:"confirmed_time"`
	CreatedAt          time.Time `json:"created_at"`
	RewardAmount       float64   `json:"amount"` // amount of reward after fee subtractions
	RewardedAmount     float64   `json:"rewarded_amount"`
	CommissionFee      float64   `json:"commission_fee"` // commission fee (in percentage)
	ImmediateRelease   float64   `json:"immediate_release"`
	RewardEndAt        time.Time `json:"reward_end_at"`
	RewardStartAt      time.Time `json:"reward_start_at"`
	BrandLogo          string    `json:"brand_logo"`
}

func BuildOrderStatusQuery(status string) (query string, params interface{}) {
	switch status {
	case dto.OrderStatusWaitForConfirming:
		query = "AND o.order_status IN ?"
		params := []string{OrderStatusInitial, OrderStatusPending, OrderStatusApproved}
		return query, params
	case dto.OrderStatusRewarding:
		return "AND o.order_status = ?", OrderStatusRewarding
	case dto.OrderStatusComplete:
		return "AND o.order_status = ?", OrderStatusComplete
	case dto.OrderStatusRejected:
		return "AND o.order_status = ?", OrderStatusRejected
	case dto.OrderStatusCancelled:
		return "AND o.order_status = ?", OrderStatusCancelled
	}
	return "", nil
}

func (o *OrderDetails) ToOrderDetailsDto() dto.OrderDetailsDto {
	var status string
	rejectedTime := time.Time{}

	switch o.OrderStatus {
	case OrderStatusRewarding:
		status = dto.OrderStatusRewarding
	case OrderStatusComplete:
		status = dto.OrderStatusComplete
	case OrderStatusRejected:
		status = dto.OrderStatusRejected
		rejectedTime = o.UpdateTime
	case OrderStatusCancelled:
		status = dto.OrderStatusCancelled
	default:
		// case OrderStatusInitial, OrderStatusPending, OrderStatusApproved
		status = dto.OrderStatusWaitForConfirming
	}

	timeToCalculateUnlockAmt := time.Now()
	if status == dto.OrderStatusRejected {
		timeToCalculateUnlockAmt = rejectedTime
	}

	daysPassed := int(timeToCalculateUnlockAmt.Sub(o.RewardStartAt) / OneDay) // number of days passed since order created
	totalDays := int(o.RewardEndAt.Sub(o.RewardStartAt) / OneDay)             // total lock days
	dayPassedPercent := float64(daysPassed) / float64(totalDays)
	if dayPassedPercent > 1 {
		dayPassedPercent = 1
	}

	imReleaseAmount := util.RoundFloat(o.RewardAmount*o.ImmediateRelease, 2)
	secondPartReward := o.RewardAmount - imReleaseAmount
	secondPartUnlockedAmount := util.RoundFloat(secondPartReward*dayPassedPercent, 2)
	rewardRemainingAmount := util.RoundFloat(secondPartReward-secondPartUnlockedAmount, 2)

	return dto.OrderDetailsDto{
		ID:                             o.ID,
		UserId:                         o.UserId,
		OrderStatus:                    status,
		Billing:                        o.Billing,
		CategoryName:                   o.CategoryName,
		Merchant:                       o.Merchant,
		ImageUrl:                       o.BrandLogo,
		AccessTradeOrderId:             o.AccessTradeOrderId,
		PubCommission:                  o.PubCommission,
		CommissionFee:                  o.CommissionFee,
		BuyTime:                        o.SalesTime,
		ConfirmedTime:                  o.CreatedAt,
		RejectedTime:                   rejectedTime,
		CancelledTime:                  o.CancelledTime,
		RewardFirstPartReleasedTime:    o.RewardStartAt,
		RewardFirstPartReleasedAmount:  imReleaseAmount,
		RewardSecondPartUnlockedAmount: secondPartUnlockedAmount,
		RewardRemainingAmount:          rewardRemainingAmount,
		RewardAmount:                   o.RewardAmount,
		RewardEndAt:                    o.RewardEndAt,
	}
}

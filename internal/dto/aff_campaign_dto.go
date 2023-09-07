package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type CampaignDescriptionDto struct {
	ID                    uint      `json:"id"`
	CampaignId            uint      `json:"campaign_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	ActionPoint           string    `json:"action_point"`
	CommissionPolicy      string    `json:"commission_policy"`
	CookiePolicy          string    `json:"cookie_policy"`
	Introduction          string    `json:"introduction"`
	OtherNotice           string    `json:"other_notice"`
	RejectedReason        string    `json:"rejected_reason"`
	TrafficBuildingPolicy string    `json:"traffic_building_policy"`
}

type AffCampaignDto struct {
	ID                uint                    `json:"id"`
	ActiveStatus      int                     `json:"active_status"`
	BrandId           uint                    `json:"brand_id"`
	AccessTradeId     string                  `json:"accesstrade_id"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
	Logo              string                  `json:"logo"`
	MaxCom            string                  `json:"max_com"`
	Merchant          string                  `json:"merchant"`
	Name              string                  `json:"name"`
	Scope             string                  `json:"scope"`
	Approval          string                  `json:"approval"`
	Status            int                     `json:"status"`
	Type              int                     `json:"type"`
	Url               string                  `json:"url"`
	Category          string                  `json:"category"`
	SubCategory       string                  `json:"sub_category"`
	CookieDuration    int                     `json:"cookie_duration"`
	CookiePolicy      string                  `json:"cookie_policy"`
	Description       *CampaignDescriptionDto `json:"description"`
	StartTime         *time.Time              `json:"start_time"`
	EndTime           *time.Time              `json:"end_time"`
	StellaDescription interface{}             `json:"stella_description"`
	CategoryId        uint                    `json:"category_id"`
	StellaStatus      string                  `json:"stella_status"`
	Thumbnail         string                  `json:"thumbnail"`
	StellaMaxCom      decimal.Decimal         `json:"stella_max_com"`
}

type AffCampaignDtoResponse struct {
	NextPage int              `json:"next_page"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
	Total    int64            `json:"total"`
	Data     []AffCampaignDto `json:"data"`
}

package dto

import (
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
	ID             uint                    `json:"id"`
	ActiveStatus   int                     `json:"active_status"`
	AccessTradeId  string                  `json:"accesstrade_id"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	Logo           string                  `json:"logo"`
	MaxCom         string                  `json:"max_com"`
	Merchant       string                  `json:"merchant"`
	Scope          string                  `json:"scope"`
	Approval       string                  `json:"approval"`
	Status         int                     `json:"status"`
	Type           int                     `json:"type"`
	Category       string                  `json:"category"`
	SubCategory    string                  `json:"sub_category"`
	CookieDuration int                     `json:"cookie_duration"`
	CookiePolicy   string                  `json:"cookie_policy"`
	Description    *CampaignDescriptionDto `json:"description"`
	StellaInfo     StellaInfoDto           `json:"stella_info"`
}

type StellaInfoDto struct {
	StartTime         *time.Time  `json:"start_time"`
	EndTime           *time.Time  `json:"end_time"`
	StellaDescription interface{} `json:"stella_description"`
	CategoryId        uint        `json:"category_id"`
	StellaStatus      string      `json:"stella_status"`
	Thumbnail         string      `json:"thumbnail"`
	StellaMaxCom      string      `json:"stella_max_com"`
	Url               string      `json:"url"`
	Name              string      `json:"name"`
	BrandId           uint        `json:"brand_id"`
	Brand             BrandDto    `json:"brand"`
	Category          CategoryDto `json:"category"`
}

type AffCampaignAppDto struct {
	ID                uint        `json:"id"`
	BrandId           uint        `json:"brand_id"`
	Brand             BrandDto    `json:"brand"`
	AccessTradeId     string      `json:"accesstrade_id"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	Name              string      `json:"name"`
	Url               string      `json:"url"`
	StartTime         *time.Time  `json:"start_time"`
	EndTime           *time.Time  `json:"end_time"`
	StellaDescription interface{} `json:"stella_description"`
	CategoryId        uint        `json:"category_id"`
	StellaStatus      string      `json:"stella_status"`
	Thumbnail         string      `json:"thumbnail"`
	StellaMaxCom      string      `json:"stella_max_com"`
}

type AffCampaignDtoResponse struct {
	NextPage int              `json:"next_page"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
	Total    int64            `json:"total"`
	Data     []AffCampaignDto `json:"data"`
}

type AffCampaignAppDtoResponse struct {
	NextPage int                 `json:"next_page"`
	Page     int                 `json:"page"`
	Size     int                 `json:"size"`
	Total    int64               `json:"total,omitempty"`
	Data     []AffCampaignAppDto `json:"data"`
}

type UserViewAffCampDto struct {
	UserId    uint32 `json:"user_id"`
	AffCampId uint64 `json:"aff_camp_id"`
}

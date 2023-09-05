package model

import (
	"time"
)

type CampaignDescription struct {
	ID                    uint      `gorm:"primarykey" json:"id"`
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

func (c *CampaignDescription) TableName() string {
	return "aff_campaign_description"
}

type Campaign struct {
	ID             uint                 `gorm:"primarykey" json:"id"`
	ActiveStatus   int                  `json:"active_status"`
	BrandId        uint                 `json:"brand_id"`
	AccesstradeId  string               `json:"accesstrade_id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Logo           string               `json:"logo"`
	MaxCom         string               `json:"max_com"`
	Merchant       string               `json:"merchant"`
	Name           string               `json:"name"`
	Scope          string               `json:"scope"`
	Approval       string               `json:"approval"`
	Status         int                  `json:"status"`
	Type           int                  `json:"type"`
	Url            string               `json:"url"`
	Category       string               `json:"category"`
	SubCategory    string               `json:"sub_category"`
	CookieDuration int                  `json:"cookie_duration"`
	CookiePolicy   string               `json:"cookie_policy"`
	Description    *CampaignDescription `json:"description" gorm:"foreignKey:CampaignId;references:ID"`
	StartTime      *time.Time           `json:"start_time"`
	EndTime        *time.Time           `json:"end_time"`
}

func (c *Campaign) TableName() string {
	return "aff_campaign"
}

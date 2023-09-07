package model

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/grokify/html-strip-tags-go"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
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

func (c *CampaignDescription) ToDto() *dto.CampaignDescriptionDto {
	return &dto.CampaignDescriptionDto{
		ID:                    c.ID,
		CampaignId:            c.CampaignId,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
		ActionPoint:           strip.StripTags(c.ActionPoint),
		CommissionPolicy:      strip.StripTags(c.CommissionPolicy),
		CookiePolicy:          strip.StripTags(c.CookiePolicy),
		Introduction:          strip.StripTags(c.Introduction),
		OtherNotice:           strip.StripTags(c.OtherNotice),
		RejectedReason:        strip.StripTags(c.RejectedReason),
		TrafficBuildingPolicy: strip.StripTags(c.TrafficBuildingPolicy),
	}
}

func (c *CampaignDescription) TableName() string {
	return "aff_campaign_description"
}

type AffCampaign struct {
	ID                uint                `gorm:"primarykey" json:"id"`
	ActiveStatus      int                 `json:"active_status"`
	BrandId           uint                `json:"brand_id"`
	AccessTradeId     string              `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Logo              string              `json:"logo"`
	MaxCom            string              `json:"max_com"`
	Merchant          string              `json:"merchant"`
	Name              string              `json:"name"`
	Scope             string              `json:"scope"`
	Approval          string              `json:"approval"`
	Status            int                 `json:"status"`
	Type              int                 `json:"type"`
	Url               string              `json:"url"`
	Category          string              `json:"category"`
	SubCategory       string              `json:"sub_category"`
	CookieDuration    int                 `json:"cookie_duration"`
	CookiePolicy      string              `json:"cookie_policy"`
	Description       CampaignDescription `json:"description" gorm:"foreignKey:CampaignId;references:ID"`
	StartTime         *time.Time          `json:"start_time"`
	EndTime           *time.Time          `json:"end_time"`
	StellaDescription datatypes.JSON      `json:"stella_description"`
	CategoryId        uint                `json:"category_id"`
	StellaStatus      string              `json:"stella_status"`
	Thumbnail         string              `json:"thumbnail"`
	StellaMaxCom      decimal.Decimal     `json:"stella_max_com" gorm:"type:decimal(4,2);"`
}

func (c *AffCampaign) TableName() string {
	return "aff_campaign"
}

func (c *AffCampaign) ToDto() dto.AffCampaignDto {
	campDto := dto.AffCampaignDto{
		ID:                c.ID,
		BrandId:           c.BrandId,
		AccessTradeId:     c.AccessTradeId,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
		Thumbnail:         c.Thumbnail,
		MaxCom:            c.MaxCom,
		Merchant:          c.Merchant,
		Name:              c.Name,
		Status:            c.Status,
		Url:               c.Url,
		CategoryId:        c.CategoryId,
		StartTime:         c.StartTime,
		EndTime:           c.EndTime,
		StellaDescription: c.StellaDescription,
		StellaStatus:      c.StellaStatus,
		StellaMaxCom:      c.StellaMaxCom,
	}
	campDto.Description = c.Description.ToDto()
	return campDto
}

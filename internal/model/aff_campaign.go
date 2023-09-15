package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	strip "github.com/grokify/html-strip-tags-go"
	"gorm.io/datatypes"
)

const (
	StellaStatusInProgress = "IN_PROGRESS"
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
	Brand             Brand               `json:"brand" gorm:"foreignKey:BrandId"`
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
	StellaCategory    Category            `json:"stella_category" gorm:"foreignKey:CategoryId;references:ID"`
	StellaStatus      string              `json:"stella_status"`
	Thumbnail         string              `json:"thumbnail"`
	StellaMaxCom      string              `json:"stella_max_com"`
}

func (c *AffCampaign) TableName() string {
	return "aff_campaign"
}

func (c *AffCampaign) ToDto() dto.AffCampaignDto {
	campDto := dto.AffCampaignDto{
		ID:            c.ID,
		AccessTradeId: c.AccessTradeId,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		MaxCom:        c.MaxCom,
		Merchant:      c.Merchant,
		Status:        c.Status,
		StellaInfo: dto.StellaInfoDto{
			Url:               c.Url,
			CategoryId:        c.CategoryId,
			StartTime:         c.StartTime,
			EndTime:           c.EndTime,
			StellaDescription: c.StellaDescription,
			StellaStatus:      c.StellaStatus,
			StellaMaxCom:      c.StellaMaxCom,
			Thumbnail:         c.Thumbnail,
			Name:              c.Name,
			BrandId:           c.BrandId,
			Brand:             c.Brand.ToBrandDto(),
			Category:          c.StellaCategory.ToCategoryDto(),
		},
	}
	campDto.Description = c.Description.ToDto()
	return campDto
}

type AffCampaignApp struct {
	ID                uint64         `gorm:"primarykey" json:"id"`
	BrandId           uint64         `json:"brand_id"`
	Brand             Brand          `json:"brand" gorm:"foreignKey:BrandId"`
	AccessTradeId     string         `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	Name              string         `json:"name"`
	Url               string         `json:"url"`
	StartTime         *time.Time     `json:"start_time"`
	EndTime           *time.Time     `json:"end_time"`
	StellaDescription datatypes.JSON `json:"stella_description"`
	CategoryId        uint64         `json:"category_id"`
	StellaStatus      string         `json:"stella_status"`
	Thumbnail         string         `json:"thumbnail"`
	StellaMaxCom      string         `json:"stella_max_com"`
}

func (c *AffCampaignApp) TableName() string {
	return "aff_campaign"
}

func (c *AffCampaignApp) ToAffCampaignAppDto() dto.AffCampaignAppDto {
	return dto.AffCampaignAppDto{
		ID:                c.ID,
		BrandId:           c.BrandId,
		Brand:             c.Brand.ToBrandDto(),
		AccessTradeId:     c.AccessTradeId,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
		Thumbnail:         c.Thumbnail,
		Name:              c.Name,
		Url:               c.Url,
		CategoryId:        c.CategoryId,
		StartTime:         c.StartTime,
		EndTime:           c.EndTime,
		StellaDescription: c.StellaDescription,
		StellaStatus:      c.StellaStatus,
		StellaMaxCom:      c.StellaMaxCom,
	}
}

type AffCampaignComBrand struct {
	ID           uint64 `gorm:"primarykey" json:"id"`
	Name         string `json:"name"`
	BrandId      uint64 `json:"brand_id"`
	Brand        Brand  `json:"brand" gorm:"foreignKey:BrandId"`
	MaxCom       string `json:"max_com"`
	StellaMaxCom string `json:"stella_max_com"`
}

func (c *AffCampaignComBrand) TableName() string {
	return "aff_campaign"
}

func (c *AffCampaignComBrand) ToAffCampaignComBrandDto() dto.AffCampaignComBrandDto {
	return dto.AffCampaignComBrandDto{
		ID:           c.ID,
		BrandId:      c.BrandId,
		Brand:        c.Brand.ToBrandDto(),
		MaxCom:       c.MaxCom,
		StellaMaxCom: c.StellaMaxCom,
	}
}

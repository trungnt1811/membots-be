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

func (m *CampaignDescription) ToDto() *dto.CampaignDescriptionDto {
	return &dto.CampaignDescriptionDto{
		ID:                    m.ID,
		CampaignId:            m.CampaignId,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		ActionPoint:           strip.StripTags(m.ActionPoint),
		CommissionPolicy:      strip.StripTags(m.CommissionPolicy),
		CookiePolicy:          strip.StripTags(m.CookiePolicy),
		Introduction:          strip.StripTags(m.Introduction),
		OtherNotice:           strip.StripTags(m.OtherNotice),
		RejectedReason:        strip.StripTags(m.RejectedReason),
		TrafficBuildingPolicy: strip.StripTags(m.TrafficBuildingPolicy),
	}
}

func (m *CampaignDescription) TableName() string {
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
	StellaCategory    AffCategory         `json:"stella_category" gorm:"foreignKey:CategoryId;references:ID"`
	StellaStatus      string              `json:"stella_status"`
	Thumbnail         string              `json:"thumbnail"`
	StellaMaxCom      string              `json:"stella_max_com"`
}

func (m *AffCampaign) TableName() string {
	return "aff_campaign"
}

func (m *AffCampaign) ToDto() dto.AffCampaignDto {
	campDto := dto.AffCampaignDto{
		ID:            m.ID,
		AccessTradeId: m.AccessTradeId,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		MaxCom:        m.MaxCom,
		Merchant:      m.Merchant,
		Status:        m.Status,
		StellaInfo: dto.StellaInfoDto{
			Url:               m.Url,
			CategoryId:        m.CategoryId,
			StartTime:         m.StartTime,
			EndTime:           m.EndTime,
			StellaDescription: m.StellaDescription,
			StellaStatus:      m.StellaStatus,
			StellaMaxCom:      m.StellaMaxCom,
			Thumbnail:         m.Thumbnail,
			Name:              m.Name,
			BrandId:           m.BrandId,
			Brand:             m.Brand.ToBrandDto(),
			Category:          m.StellaCategory.ToDto(),
		},
	}
	campDto.Description = m.Description.ToDto()
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

func (m *AffCampaignApp) TableName() string {
	return "aff_campaign"
}

func (m *AffCampaignApp) ToAffCampaignAppDto() dto.AffCampaignAppDto {
	return dto.AffCampaignAppDto{
		ID:                m.ID,
		BrandId:           m.BrandId,
		Brand:             m.Brand.ToBrandDto(),
		AccessTradeId:     m.AccessTradeId,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		Thumbnail:         m.Thumbnail,
		Name:              m.Name,
		Url:               m.Url,
		CategoryId:        m.CategoryId,
		StartTime:         m.StartTime,
		EndTime:           m.EndTime,
		StellaDescription: m.StellaDescription,
		StellaStatus:      m.StellaStatus,
		StellaMaxCom:      m.StellaMaxCom,
	}
}

type AffCampaignComBrand struct {
	ID           uint64 `gorm:"primarykey" json:"id"`
	Name         string `json:"name"`
	BrandId      uint64 `json:"brand_id"`
	Brand        Brand  `json:"brand" gorm:"foreignKey:BrandId"`
	StellaMaxCom string `json:"stella_max_com"`
}

func (m *AffCampaignComBrand) TableName() string {
	return "aff_campaign"
}

func (m *AffCampaignComBrand) ToAffCampaignLessDto() dto.AffCampaignLessDto {
	return dto.AffCampaignLessDto{
		ID:           uint(m.ID),
		Name:         m.Name,
		BrandId:      m.BrandId,
		Brand:        m.Brand.ToBrandDto(),
		StellaMaxCom: m.StellaMaxCom,
	}
}

type AffCampComFavBrand struct {
	ID            uint64            `gorm:"primarykey" json:"id"`
	Name          string            `json:"name"`
	BrandId       uint64            `json:"brand_id"`
	FavoriteBrand UserFavoriteBrand `gorm:"foreignKey:BrandId;references:BrandId" json:"favrorite_brand"`
	StellaMaxCom  string            `json:"stella_max_com"`
}

func (m *AffCampComFavBrand) TableName() string {
	return "aff_campaign"
}

func (m *AffCampComFavBrand) ToAffCampaignLessDto() dto.AffCampaignLessDto {
	return dto.AffCampaignLessDto{
		ID:           uint(m.ID),
		Name:         m.Name,
		BrandId:      m.BrandId,
		Brand:        m.FavoriteBrand.Brand.ToBrandDto(),
		StellaMaxCom: m.StellaMaxCom,
	}
}

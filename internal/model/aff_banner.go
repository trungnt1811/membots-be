package model

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"time"
)

func (e *AffBanner) TableName() string {
	return "aff_banner"
}

type AffBanner struct {
	ID            uint32          `json:"id" gorm:"primaryKey"`
	Name          string          `json:"name"`
	Thumbnail     string          `json:"thumbnail"`
	Url           string          `json:"url"`
	AccessTradeId string          `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	AffCampaign   AffCampaignLess `json:"aff_campaign" gorm:"foreignKey:AccessTradeId;references:AccessTradeId"`
	CreatedAt     *time.Time      `json:"created_at"`
	Status        string          `json:"status"`
}

func (c *AffBanner) ToDto() dto.AffBannerDto {
	return dto.AffBannerDto{
		ID:            c.ID,
		Name:          c.Name,
		Thumbnail:     c.Thumbnail,
		Url:           c.Url,
		AccessTradeId: c.AccessTradeId,
		CreatedAt:     c.CreatedAt,
		Status:        c.Status,
		AffCampaign:   c.AffCampaign.ToDto(),
	}
}

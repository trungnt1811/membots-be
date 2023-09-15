package model

import "github.com/astraprotocol/affiliate-system/internal/dto"

type AffCampaignLess struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	AccessTradeId string `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	Name          string `json:"name"`
	Url           string `json:"url"`
}

func (c *AffCampaignLess) TableName() string {
	return "aff_campaign"
}

func (c *AffCampaignLess) ToDto() dto.AffCampaignLessDto {
	return dto.AffCampaignLessDto{
		ID:            c.ID,
		Name:          c.Name,
		AccessTradeId: c.AccessTradeId,
		Url:           c.Url,
	}
}

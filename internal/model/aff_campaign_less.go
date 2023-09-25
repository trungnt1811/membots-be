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

type AffCampaignLessApp struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	AccessTradeId string `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	BrandId       uint   `json:"brand_id"`
	StellaStatus  string `json:"stella_status"`
	Brand         Brand  `json:"brand" gorm:"foreignKey:BrandId"`
	CategoryId    uint   `json:"category_id"`
	StellaMaxCom  string `json:"stella_max_com"`
}

func (c *AffCampaignLessApp) TableName() string {
	return "aff_campaign"
}

func (c *AffCampaignLessApp) ToDto() dto.AffCampaignLessDto {
	return dto.AffCampaignLessDto{
		ID:            c.ID,
		Name:          c.Name,
		AccessTradeId: c.AccessTradeId,
		Url:           c.Url,
		BrandId:       c.BrandId,
		Brand:         c.Brand.ToBrandDto(),
		StellaStatus:  c.StellaStatus,
		StellaMaxCom:  c.StellaMaxCom,
	}
}

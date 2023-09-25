package model

import "github.com/astraprotocol/affiliate-system/internal/dto"

type AffCampaignLess struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	AccessTradeId string `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	Name          string `json:"name"`
	Url           string `json:"url"`
}

func (m *AffCampaignLess) TableName() string {
	return "aff_campaign"
}

func (m *AffCampaignLess) ToDto() dto.AffCampaignLessDto {
	return dto.AffCampaignLessDto{
		ID:            m.ID,
		Name:          m.Name,
		AccessTradeId: m.AccessTradeId,
		Url:           m.Url,
	}
}

type AffCampaignLessApp struct {
	ID            uint                   `gorm:"primarykey" json:"id"`
	AccessTradeId string                 `json:"accesstrade_id" gorm:"column:accesstrade_id"`
	Name          string                 `json:"name"`
	Url           string                 `json:"url"`
	BrandId       uint64                 `json:"brand_id"`
	StellaStatus  string                 `json:"stella_status"`
	Brand         Brand                  `json:"brand" gorm:"foreignKey:BrandId"`
	CategoryId    uint                   `json:"category_id"`
	StellaMaxCom  string                 `json:"stella_max_com"`
	Attributes    []AffCampaignAttribute `json:"attributes" gorm:"foreignKey:CampaignId"`
}

func (m *AffCampaignLessApp) TableName() string {
	return "aff_campaign"
}

func (m *AffCampaignLessApp) ToDto() dto.AffCampaignLessDto {
	var listAttribute []dto.AffCampaignAttributeDto
	for _, attribute := range m.Attributes {
		listAttribute = append(listAttribute, attribute.ToDto())
	}
	return dto.AffCampaignLessDto{
		ID:            m.ID,
		Name:          m.Name,
		AccessTradeId: m.AccessTradeId,
		Url:           m.Url,
		BrandId:       m.BrandId,
		Brand:         m.Brand.ToBrandDto(),
		StellaStatus:  m.StellaStatus,
		StellaMaxCom:  m.StellaMaxCom,
		Attributes:    listAttribute,
	}
}

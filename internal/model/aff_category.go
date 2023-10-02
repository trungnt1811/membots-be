package model

import "github.com/astraprotocol/affiliate-system/internal/dto"

func (e *AffCategory) TableName() string {
	return "category"
}

type AffCategory struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

func (c *AffCategory) ToDto() dto.AffCategoryDto {
	return dto.AffCategoryDto{
		ID:   c.ID,
		Logo: c.Logo,
		Name: c.Name,
	}
}

type AffCategoryAndTotalCampaign struct {
	ID               uint64 `json:"id" gorm:"primaryKey"`
	Name             string `json:"name"`
	Logo             string `json:"logo"`
	TotalAffCampaign uint32 `json:"total_aff_campaign"`
}

func (c *AffCategoryAndTotalCampaign) ToDto() dto.AffCategoryDto {
	return dto.AffCategoryDto{
		ID:               c.ID,
		Logo:             c.Logo,
		Name:             c.Name,
		TotalAffCampaign: c.TotalAffCampaign,
	}
}

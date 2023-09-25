package model

import "github.com/astraprotocol/affiliate-system/internal/dto"

var AttributeTypePriorityMapping = map[string]int{
	"percent": 1,
	"vnd":     2,
}

type AffCampaignAttribute struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	CampaignId     uint   `json:"campaign_id"`
	AttributeKey   string `json:"attribute_key"`
	AttributeValue string `json:"attribute_value"`
	AttributeType  string `json:"attribute_type"`
}

func (a *AffCampaignAttribute) TableName() string {
	return "aff_campaign_attribute"
}

func (a *AffCampaignAttribute) ToDto() dto.AffCampaignAttributeDto {
	return dto.AffCampaignAttributeDto{
		ID:             a.ID,
		CampaignId:     a.CampaignId,
		AttributeKey:   a.AttributeKey,
		AttributeValue: a.AttributeValue,
		AttributeType:  a.AttributeType,
	}
}

package model

const SearchHistoryTypeAffCampaign = "aff_campaign"
const TrendTypeAffCampaign = "aff_campaign"

type AffSearch struct {
	AffCampaign   []AffCampaignLessApp `json:"aff_campaign"`
	TotalCampaign int64                `json:"total_campaign"`
}

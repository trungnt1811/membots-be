package dto

type AffSearchDto struct {
	AffCampaigns  []AffCampaignAppDto `json:"aff_campaigns"`
	TotalCampaign int64               `json:"total_campaign"`
}

type AffSearchResponseDto struct {
	NextPage int          `json:"next_page"`
	Page     int          `json:"page"`
	Size     int          `json:"size"`
	Data     AffSearchDto `json:"data"`
}

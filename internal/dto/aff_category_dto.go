package dto

type AffCategoryDto struct {
	ID               uint64 `json:"id"`
	Name             string `json:"name"`
	Logo             string `json:"logo"`
	StellaMaxCom     string `json:"stella_max_com,omitempty"`
	TotalAffCampaign uint32 `json:"total_aff_campaign,omitempty"`
}
type AffCategoryResponseDto struct {
	NextPage int              `json:"next_page"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
	Data     []AffCategoryDto `json:"data"`
}

package dto

type AffCampaignLessDto struct {
	ID            uint     `json:"campaign_id"`
	AccessTradeId string   `json:"accesstrade_id"`
	Name          string   `json:"name"`
	Url           string   `json:"url"`
	StellaStatus  string   `json:"stella_status,omitempty"`
	BrandId       uint64   `json:"brand_id,omitempty"`
	Brand         BrandDto `json:"brand,omitempty"`
}

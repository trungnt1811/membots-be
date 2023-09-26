package dto

type AffCampaignLessDto struct {
	ID            uint     `json:"campaign_id"`
	AccessTradeId string   `json:"accesstrade_id,omitempty"`
	Name          string   `json:"name"`
	Url           string   `json:"url,omitempty"`
	StellaStatus  string   `json:"stella_status,omitempty"`
	StellaMaxCom  string   `json:"stella_max_com,omitempty"`
	BrandId       uint     `json:"brand_id,omitempty"`
	Brand         BrandDto `json:"brand,omitempty"`
}

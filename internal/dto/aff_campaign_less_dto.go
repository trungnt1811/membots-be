package dto

type AffCampaignLessDto struct {
	ID            uint   `json:"campaign_id"`
	AccessTradeId string `json:"accesstrade_id"`
	Name          string `json:"name"`
	Url           string `json:"url"`
}

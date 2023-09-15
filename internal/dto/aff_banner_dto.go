package dto

import "time"

type AffBannerDto struct {
	ID            uint32             `json:"id"`
	Name          string             `json:"name"`
	Thumbnail     string             `json:"thumbnail"`
	Url           string             `json:"url"`
	AccessTradeId string             `json:"accesstrade_id"`
	CreatedAt     *time.Time         `json:"created_at"`
	Status        string             `json:"status"`
	AffCampaign   AffCampaignLessDto `json:"aff_campaign"`
}

type AffBannerDtoResponse struct {
	NextPage int            `json:"next_page"`
	Page     int            `json:"page"`
	Size     int            `json:"size"`
	Total    int64          `json:"total"`
	Data     []AffBannerDto `json:"data"`
}

type AffBannerCreateDto struct {
	Name          string `json:"name"`
	Thumbnail     string `json:"thumbnail"`
	Url           string `json:"url"`
	AccessTradeId string `json:"accesstrade_id"`
}

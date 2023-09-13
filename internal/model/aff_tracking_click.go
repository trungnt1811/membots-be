package model

import "time"

type AffTrackedClick struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CampaignId uint      `json:"campaign_id"`
	UserId     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	AffLink    string    `json:"aff_link"`
	ShortLink  string    `json:"short_link"`
	UrlOrigin  string    `json:"url_origin"`
	OrderId    string    `json:"order_id"`
}

func (e *AffTrackedClick) TableName() string {
	return "aff_tracked_click"
}

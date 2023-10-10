package model

import "time"

type AffTrackedClick struct {
	ID         uint64       `gorm:"primarykey" json:"id"`
	CampaignId uint         `json:"campaign_id"`
	UserId     uint         `json:"user_id"`
	LinkId     uint         `json:"link_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	AffLink    string       `json:"aff_link"`
	UrlOrigin  string       `json:"url_origin"`
	OrderId    string       `json:"order_id"`
	Campaign   *AffCampaign `json:"campaign" gorm:"foreignKey:CampaignId"`
}

func (e *AffTrackedClick) TableName() string {
	return "aff_tracked_click"
}

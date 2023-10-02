package model

import "time"

const (
	AFF_LINK_STATUS_INACTIVE = 0
	AFF_LINK_STATUS_ACTIVE   = 1
)

type AffLink struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	AffLink      string    `json:"aff_link"`
	FirstLink    string    `json:"first_link"`
	ShortLink    string    `json:"short_link"`
	UrlOrigin    string    `json:"url_origin"`
	CampaignId   uint      `json:"campaign_id"`
	ActiveStatus int8      `json:"active_status"`
}

func (e *AffLink) TableName() string {
	return "aff_link"
}

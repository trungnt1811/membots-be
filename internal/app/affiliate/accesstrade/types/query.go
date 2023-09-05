package types

import "time"

const (
	ORDER_STATUS_PENDING  = 0
	ORDER_STATUS_APPROVED = 1
	ORDER_STATUS_REJECTED = 2
)

type ATOrderQuery struct {
	Since       time.Time `url:"since"`
	Until       time.Time `url:"until"`
	Status      int       `url:"status,omitempty"`
	Merchant    string    `url:"merchant,omitempty"`
	UTMSource   string    `url:"utm_source,omitempty"`
	UTMCampaign string    `url:"utm_campaign,omitempty"`
	UTMMedium   string    `url:"utm_medium,omitempty"`
	UTMContent  string    `url:"utm_content,omitempty"`
}

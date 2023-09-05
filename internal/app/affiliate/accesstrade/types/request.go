package types

type ATLinkReq struct {
	CampaignId  string   `json:"campaign_id"`
	Urls        []string `json:"urls,omitempty"`
	UrlEnc      bool     `json:"url_enc,omitempty"`
	UTMSource   string   `json:"utm_source,omitempty"`
	UTMMedium   bool     `json:"utm_medium,omitempty"`
	UTMCampaign bool     `json:"utm_campaign,omitempty"`
	UTMContent  bool     `json:"utm_content,omitempty"`
	Sub1        bool     `json:"sub1,omitempty"`
	Sub2        bool     `json:"sub2,omitempty"`
	Sub3        bool     `json:"sub3,omitempty"`
	Sub4        bool     `json:"sub4,omitempty"`
}

package dto

import "errors"

type CreateLinkPayload struct {
	CampaignId  uint   `json:"campaign_id"`
	OriginalUrl string `json:"original_url"`
	ShortenLink bool   `json:"shorten_link"`
}

func (p *CreateLinkPayload) IsValid() error {
	if p.CampaignId == 0 {
		return errors.New("invalid campaign id")
	}

	return nil
}

type CreateLinkResponse struct {
	CampaignId  uint   `json:"campaign_id"`
	AffLink     string `json:"aff_link"`
	ShortLink   string `json:"short_link"`
	OriginalUrl string `json:"original_url"`
}

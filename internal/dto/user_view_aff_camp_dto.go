package dto

import "time"

type UserViewAffCampComBrandDto struct {
	ID              uint64                 `json:"id"`
	UserId          uint32                 `json:"user_id"`
	AffCampId       uint64                 `json:"aff_camp_id"`
	AffCampComBrand AffCampaignComBrandDto `json:"aff_campaign"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

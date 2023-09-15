package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

type UserViewAffCamp struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserId    uint32    `json:"user_id"`
	AffCampId uint64    `json:"aff_camp_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *UserViewAffCamp) TableName() string {
	return "user_view_aff_camp"
}

type UserViewAffCampComBrand struct {
	AffCampId       uint64              `json:"aff_camp_id"`
	AffCampComBrand AffCampaignComBrand `json:"aff_campaign" gorm:"foreignKey:AffCampId"`
}

func (m *UserViewAffCampComBrand) TableName() string {
	return "user_view_aff_camp"
}

func (m *UserViewAffCampComBrand) ToAffCampaignComBrandDto() dto.AffCampaignComBrandDto {
	return m.AffCampComBrand.ToAffCampaignComBrandDto()
}

package model

import (
	"time"
)

func (m *UserViewBrand) TableName() string {
	return "user_view_brand"
}

type UserViewBrand struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	UserId    uint32    `json:"user_id"`
	AffCampId uint64    `json:"aff_camp_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

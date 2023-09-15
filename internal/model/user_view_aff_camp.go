package model

import (
	"time"
)

func (m *UserViewAffCamp) TableName() string {
	return "user_view_aff_camp"
}

type UserViewAffCamp struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	UserId    uint32    `json:"user_id"`
	AffCampId uint64    `json:"aff_camp_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

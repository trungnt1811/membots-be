package model

import (
	"time"
)

const (
	UserFavoriteBrandStatusAdded = "ADDED"
)

type UserFavoriteBrand struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserId    uint32    `json:"user_id"`
	BrandId   uint64    `json:"brand_id"`
	Brand     Brand     `json:"brand" gorm:"foreignKey:BrandId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *UserFavoriteBrand) TableName() string {
	return "user_favorite_brand"
}

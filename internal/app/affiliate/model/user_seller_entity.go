package model

import "time"

type UserSellerEntity struct {
	UserId    int       `json:"userId" gorm:"primaryKey"`
	SellerId  int       `json:"sellerId" gorm:"primaryKey"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *UserSellerEntity) TableName() string {
	return "user_seller"
}

package model

import "time"

type AffBrand struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *AffBrand) TableName() string {
	return "brand"
}

type AffMerchantBrand struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Merchant string    `json:"merchant"`
	BrandId  uint      `json:"brand_id"`
	Brand    *AffBrand `gorm:"foreignKey:ID;references:BrandId" json:"brand"`
}

func (c *AffMerchantBrand) TableName() string {
	return "aff_merchant_brand"
}

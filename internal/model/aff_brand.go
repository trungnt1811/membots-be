package model

import "github.com/astraprotocol/affiliate-system/internal/dto"

func (e *Brand) TableName() string {
	return "brand"
}

type Brand struct {
	ID         uint32  `json:"id" gorm:"primaryKey"`
	Name       string  `json:"name"`
	Logo       string  `json:"logo"`
	CoverPhoto *string `json:"cover_photo"`
}

func (c *Brand) ToBrandDto() dto.BrandDto {
	return dto.BrandDto{
		ID:         c.ID,
		Logo:       c.Logo,
		Name:       c.Name,
		CoverPhoto: c.CoverPhoto,
	}
}

type AffMerchantBrand struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Merchant string `json:"merchant"`
	BrandId  uint   `json:"brand_id"`
	Brand    *Brand `gorm:"foreignKey:ID;references:BrandId" json:"brand"`
}

func (c *AffMerchantBrand) TableName() string {
	return "aff_merchant_brand"
}

type TotalFavoriteBrand struct {
	BrandId       uint64 `json:"brand_id"`
	TotalFavorite uint64 `json:"total_fav"`
}

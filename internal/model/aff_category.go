package model

import "github.com/astraprotocol/affiliate-system/internal/dto"

func (e *Category) TableName() string {
	return "category"
}

type Category struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

func (c *Category) ToCategoryDto() dto.CategoryDto {
	return dto.CategoryDto{
		ID:   c.ID,
		Logo: c.Logo,
		Name: c.Name,
	}
}

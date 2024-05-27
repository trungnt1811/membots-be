package model

import "github.com/flexstack.ai/membots-be/internal/dto"

type Meme20Meta struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	Ticker     string  `json:"ticker"`
}

func (e *Meme20Meta) TableName() string {
	return "meme20_meta"
}

func (c *Meme20Meta) ToDto() dto.Meme20MetaDto {
	return dto.Meme20MetaDto{
		ID:         c.ID,
		Ticker:     c.Ticker,
	}
}

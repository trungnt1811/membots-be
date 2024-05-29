package memeception

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
	"gorm.io/gorm"
)

type memeceptionRepository struct {
	db *gorm.DB
}

func NewMemeceptionRepository(db *gorm.DB) interfaces.MemeceptionRepository {
	return &memeceptionRepository{
		db: db,
	}
}

func (r memeceptionRepository) GetMemeceptionBySymbol(ctx context.Context, symbol string) (model.Meme, error) {
	var memeMeta model.Meme
	err := r.db.Where("symbol = ?", symbol).Preload("Memeception").Preload("Social").First(&memeMeta).Error
	return memeMeta, err
}

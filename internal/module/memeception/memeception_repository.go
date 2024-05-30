package memeception

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
)

const beforeLauchStart = 30 // mins

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
	err := r.db.Joins("Memeception").Joins("Social").
		Where("symbol = ?", symbol).
		First(&memeMeta).Error
	return memeMeta, err
}

func (r memeceptionRepository) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("start_at < ?", time.Now().Unix()).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsUpcoming(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("start_at > ?", time.Now().Add(beforeLauchStart*time.Minute).Unix()).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where(
			"ama = ? AND start_at >= ? AND start_at <= ?",
			true,
			time.Now().Unix(),
			time.Now().Add(beforeLauchStart*time.Minute).Unix(),
		).
		Find(&memeceptions).Error
	return memeceptions, err
}

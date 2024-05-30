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
	err := r.db.Where("symbol = ?", symbol).Preload("Memeception").Preload("Social").First(&memeMeta).Error
	return memeMeta, err
}

func (r memeceptionRepository) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Where("start_at < ?", time.Now().Unix()).Preload("Meme").Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsUpcoming(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Where("start_at > ?", time.Now().Add(beforeLauchStart*time.Minute).Unix()).
		Preload("Meme").Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Where(
		"ama = ? AND start_at >= ? AND start_at <= ?",
		true,
		time.Now(),
		time.Now().Add(beforeLauchStart*time.Minute).Unix(),
	).Preload("Meme").Find(&memeceptions).Error
	return memeceptions, err
}

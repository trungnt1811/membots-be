package memeception

import (
	"context"

	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type status uint

const (
	PROCESSING status = 0
	SUCCEED    status = 1
)

type memeceptionRepository struct {
	db *gorm.DB
}

func NewMemeceptionRepository(db *gorm.DB) interfaces.MemeceptionRepository {
	return &memeceptionRepository{
		db: db,
	}
}

func (c memeceptionRepository) CreateMeme(ctx context.Context, model model.Meme) error {
	// TODO: implement later
	return nil
}

func (r memeceptionRepository) GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error) {
	var memeMeta model.Meme
	err := r.db.Joins("Memeception").Joins("Social").
		Where("contract_address = ?", contractAddress).
		Where("status = ?", SUCCEED).
		First(&memeMeta).Error
	return memeMeta, err
}

func (r memeceptionRepository) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("target_eth <= collected_eth").
		Where("status = ?", SUCCEED).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("target_eth > collected_eth").
		Where("status = ?", SUCCEED).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("status = ?", SUCCEED).
		Order("id DESC").
		Find(&memeceptions).Error
	return memeceptions, err
}

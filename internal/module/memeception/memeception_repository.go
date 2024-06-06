package memeception

import (
	"context"

	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/internal/constant"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type memeceptionRepository struct {
	db *gorm.DB
}

func NewMemeceptionRepository(db *gorm.DB) interfaces.MemeceptionRepository {
	return &memeceptionRepository{
		db: db,
	}
}

func (r memeceptionRepository) CreateMeme(ctx context.Context, model model.Meme) error {
	return r.db.Create(&model).Error
}

func (r memeceptionRepository) UpdateMeme(ctx context.Context, model model.Meme) error {
	return r.db.Updates(&model).Error
}

func (r memeceptionRepository) GetListMemeProcessing(ctx context.Context) ([]model.MemeOnchainInfo, error) {
	var meme []model.MemeOnchainInfo
	err := r.db.Where("status = ?", constant.PROCESSING).
		Find(&meme).Error
	return meme, err
}

func (r memeceptionRepository) GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error) {
	var meme model.Meme
	err := r.db.Joins("Memeception").Joins("Social").
		Where("contract_address = ?", contractAddress).
		Where("Meme.status = ?", constant.SUCCEED).
		First(&meme).Error
	return meme, err
}

func (r memeceptionRepository) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("target_eth <= collected_eth").
		Where("Meme.status = ?", constant.SUCCEED).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("target_eth > collected_eth").
		Where("Meme.status = ?", constant.SUCCEED).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("status = ?", constant.SUCCEED).
		Order("id DESC").
		Find(&memeceptions).Error
	return memeceptions, err
}

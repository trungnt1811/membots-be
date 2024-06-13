package memeception

import (
	"context"
	"errors"

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
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&model).Error
}

func (r memeceptionRepository) UpdateMemeception(ctx context.Context, model model.Memeception) error {
	return r.db.Updates(&model).Error
}

func (r memeceptionRepository) GetListMemeProcessing(ctx context.Context) ([]model.MemeOnchainInfo, error) {
	var meme []model.MemeOnchainInfo
	err := r.db.Joins("Memeception").Where("meme.status = ?", constant.PROCESSING).
		Find(&meme).Error
	return meme, err
}

func (r memeceptionRepository) GetListMemeLive(ctx context.Context) ([]model.MemeOnchainInfo, error) {
	var meme []model.MemeOnchainInfo
	err := r.db.Joins("Memeception").Where("meme.status = ?", constant.SUCCEED).
		Where("Memeception.status = ?", constant.LIVE).
		Find(&meme).Error
	return meme, err
}

func (r memeceptionRepository) GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error) {
	var meme model.Meme
	err := r.db.Joins("Memeception").Joins("Social").
		Where("meme.contract_address = ?", contractAddress).
		Where("meme.status = ?", constant.SUCCEED).
		Last(&meme).Error
	return meme, err
}

func (r memeceptionRepository) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("Meme.status = ?", constant.SUCCEED).
		Where("memeception.status = ?", constant.ENDED_SOLD_OUT).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("Meme.status = ?", constant.SUCCEED).
		Where("memeception.status = ?", constant.LIVE).
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error) {
	var memeceptions []model.Memeception
	err := r.db.Joins("Meme").
		Where("Meme.status = ?", constant.SUCCEED).
		Order("id DESC").
		Find(&memeceptions).Error
	return memeceptions, err
}

func (r memeceptionRepository) GetMapMemeSymbolAndLogoURL(ctx context.Context, contractAddresses []string) (map[string]model.MemeSymbolAndLogoURL, error) {
	var memes []model.MemeSymbolAndLogoURL
	err := r.db.Where("status = ?", constant.SUCCEED).
		Where("contract_address IN (?)", contractAddresses).
		Find(&memes).Error
	if err != nil {
		return nil, err
	}
	memeMap := make(map[string]model.MemeSymbolAndLogoURL)
	for _, meme := range memes {
		memeMap[meme.ContractAddress] = meme
	}
	return memeMap, err
}

func (r memeceptionRepository) GetMemeceptionBySymbol(ctx context.Context, symbol string) (model.Meme, error) {
	var meme model.Meme
	err := r.db.Joins("Memeception").Joins("Social").
		Where("meme.symbol = ?", symbol).
		Last(&meme).Error
	return meme, err
}

func (r memeceptionRepository) GetMemeIDAndStartAtByContractAddress(ctx context.Context, contractAddress string) (
	model.MemeceptionMemeIDAndStartAt,
	error,
) {
	var meme model.MemeceptionMemeIDAndStartAt
	err := r.db.Joins("Meme").
		Where("Meme.contract_address = ?", contractAddress).
		Last(&meme).Error
	return meme, err
}

func (r memeceptionRepository) MemeceptionExists(ctx context.Context, symbol string) (bool, error) {
	var meme model.Meme
	err := r.db.Where("meme.symbol = ?", symbol).
		Where("meme.status = ?", constant.SUCCEED).
		First(&meme).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Record not found, return false and no error
		}
		return false, err
	}
	return true, nil
}

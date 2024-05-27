package fair_launch

import (
	"context"

	"github.com/astraprotocol/membots-be/internal/interfaces"
	"github.com/astraprotocol/membots-be/internal/model"
	"gorm.io/gorm"
)

type fairLauchRepository struct {
	db *gorm.DB
}

func NewFairLauchRepository(db *gorm.DB) interfaces.FairLauchRepository {
	return &fairLauchRepository{
		db: db,
	}
}

func (r fairLauchRepository) GetMeme20MetaByTicker(ctx context.Context, ticker string) (model.Meme20Meta, error) {
	query := ""
	var meme20Meta model.Meme20Meta
	err := r.db.Raw(query, ticker).Scan(&meme20Meta).Error
	return meme20Meta, err
}

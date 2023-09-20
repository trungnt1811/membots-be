package user_favorite_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type userFavoriteBrandRepository struct {
	db *gorm.DB
}

func NewUserFavoriteBrandRepository(db *gorm.DB) interfaces.UserFavoriteBrandRepository {
	return &userFavoriteBrandRepository{
		db: db,
	}
}

func (r userFavoriteBrandRepository) GetListFavBrandByUserIdAndBrandIds(ctx context.Context, userId int64, brandIds []uint64) ([]model.UserFavoriteBrand, error) {
	var listUserFavoriteBrand []model.UserFavoriteBrand
	err := r.db.Where("user_id = ? AND status = ? AND brand_id IN ?", userId, model.UserFavoriteBrandStatusAdded, brandIds).
		Find(&listUserFavoriteBrand).Error
	return listUserFavoriteBrand, err
}

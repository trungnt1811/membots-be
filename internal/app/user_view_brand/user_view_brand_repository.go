package user_view_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type userViewBrandRepository struct {
	db *gorm.DB
}

func (r userViewBrandRepository) CreateUserViewBrand(ctx context.Context, data *model.UserViewBrand) error {
	query := "INSERT INTO user_view_brand (user_id, brand_id) " +
		"VALUES (?, ?) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP"
	return r.db.Exec(query, data.UserId, data.BrandId).Error
}

func NewUserViewBrandRepository(db *gorm.DB) interfaces.UserViewBrandRepository {
	return &userViewBrandRepository{
		db: db,
	}
}

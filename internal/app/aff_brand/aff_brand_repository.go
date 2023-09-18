package aff_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affBrandRepository struct {
	db *gorm.DB
}

func (r affBrandRepository) GetListCountFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error) {
	query := "SELECT ufb.brand_id, COUNT(*) AS total_fav FROM user_favorite_brand AS ufb " +
		"LEFT JOIN aff_campaign AS ac ON ufb.brand_id = ac.brand_id " +
		"WHERE ufb.status = 'ADDED' AND ac.brand_id != 0 " +
		"GROUP BY ufb.brand_id ORDER BY total_fav DESC"
	var listFavouriteAffBrand []model.TotalFavoriteBrand
	err := r.db.Raw(query).Row().Scan(&listFavouriteAffBrand)
	return listFavouriteAffBrand, err
}

func (r affBrandRepository) UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error {
	// must be implemented at cache layer
	return nil
}

func NewAffBrandRepository(db *gorm.DB) interfaces.AffBrandRepository {
	return &affBrandRepository{
		db: db,
	}
}

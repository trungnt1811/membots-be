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
	query := "SELECT DISTINCT ac.brand_id, cf.total_fav " +
		"FROM aff_campaign AS ac " +
		"LEFT JOIN ( " +
		"SELECT ufb.brand_id, COUNT(ufb.brand_id) AS total_fav FROM user_favorite_brand AS ufb " +
		"WHERE ufb.status = 'ADDED' " +
		"GROUP BY ufb.brand_id) AS cf " +
		"ON cf.brand_id = ac.brand_id " +
		"WHERE ac.brand_id != 0 " +
		"ORDER BY cf.total_fav DESC"
	var listFavouriteAffBrand []model.TotalFavoriteBrand
	err := r.db.Raw(query).Scan(&listFavouriteAffBrand).Error
	return listFavouriteAffBrand, err
}

func (r affBrandRepository) UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error {
	// must be implemented at cache layer
	return nil
}

func (r affBrandRepository) GetListFavAffBrandByUserId(ctx context.Context, userId uint64, page, size int) ([]model.AffCampComFavBrand, error) {
	var listAffCampComFavBrand []model.AffCampComFavBrand
	offset := (page - 1) * size
	err := r.db.Joins("JOIN FavoriteBrand ON FavoriteBrand.BrandId = aff_campaign.brand_id").
		Joins("FavoriteBrand.Brand").
		Where("user_id = ? AND FavoriteBrand.Status = ? AND aff_campaign.stella_status = ?",
			userId,
			model.UserFavoriteBrandStatusAdded,
			model.StellaStatusInProgress,
		).
		Limit(size + 1).
		Offset(offset).
		Order("FavoriteBrand.UpdatedAt DESC").
		Find(&listAffCampComFavBrand).Error
	return listAffCampComFavBrand, err
}

func NewAffBrandRepository(db *gorm.DB) interfaces.AffBrandRepository {
	return &affBrandRepository{
		db: db,
	}
}

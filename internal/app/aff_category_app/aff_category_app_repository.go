package category

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type affCategoryRepository struct {
	Db *gorm.DB
}

func (a *affCategoryRepository) GetAllCategory(ctx context.Context, page, size int) ([]model.AffCategoryAndTotalCampaign, error) {
	var listCategory []model.AffCategoryAndTotalCampaign
	offset := (page - 1) * size
	query := "select c.logo, c.id, c.name, count(c.id) as total_aff_campaign from aff_campaign as ac " +
		"left join category as c ON ac.category_id = c.id " +
		"where ac.status = ? GROUP BY c.id ORDER BY c.rank asc LIMIT ? OFFSET ?"
	err := a.Db.Raw(query, "IN_PROGRESS", size+1, offset).Scan(&listCategory).Error
	return listCategory, err
}

func (a *affCategoryRepository) GetAllAffCampaignInCategory(ctx context.Context, categoryId uint32, orderBy string, page, size int) ([]model.AffCampaignLessApp, error) {
	offset := (page - 1) * size
	var listCouponInCategory []model.AffCampaignLessApp
	err := a.Db.Table("aff_campaign").
		Where("category_id = ? AND status = ?", categoryId, "IN_PROGRESS").
		Limit(size).
		Offset(offset).
		Order(orderBy).Error
	return listCouponInCategory, err
}

func NewAppCategoryRepository(db *gorm.DB) interfaces.AffCategoryRepository {
	return &affCategoryRepository{
		Db: db,
	}
}

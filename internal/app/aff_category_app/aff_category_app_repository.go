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

func (a *affCategoryRepository) GetAttributeInCategories(ctx context.Context, ids []uint64) ([]model.CategoryWithCommissionAttribute, error) {
	var listCategory []model.CategoryWithCommissionAttribute
	query := "select aca.attribute_value, c.id, aca.attribute_type " +
		"from aff_campaign as ac left join category as c ON ac.category_id = c.id " +
		"left join (SELECT campaign_id, attribute_value, attribute_type " +
		"	FROM aff_campaign_attribute " +
		"	ORDER BY CASE WHEN attribute_type = 'percent' THEN 3 " +
		"	WHEN attribute_type = 'vnd' THEN 2 ELSE 1 END DESC, CAST(attribute_value + 0 AS DECIMAL(12,2)) DESC) as aca " +
		"	ON aca.campaign_id = ac.id " +
		"where c.id IN ? AND ac.stella_status = ?"
	err := a.Db.Raw(query, ids, "IN_PROGRESS").Scan(&listCategory).Error
	return listCategory, err
}

func (a *affCategoryRepository) GetAllCategory(ctx context.Context, page, size int) ([]model.AffCategoryAndTotalCampaign, error) {
	var listCategory []model.AffCategoryAndTotalCampaign
	offset := (page - 1) * size
	query := "select c.logo, c.id, c.name, count(c.id) as total_aff_campaign from aff_campaign as ac " +
		"left join category as c ON ac.category_id = c.id " +
		"where ac.stella_status = ? GROUP BY c.id ORDER BY c.rank asc LIMIT ? OFFSET ?"
	err := a.Db.Raw(query, "IN_PROGRESS", size+1, offset).Scan(&listCategory).Error
	return listCategory, err
}

func NewAppCategoryRepository(db *gorm.DB) interfaces.AffCategoryRepository {
	return &affCategoryRepository{
		Db: db,
	}
}

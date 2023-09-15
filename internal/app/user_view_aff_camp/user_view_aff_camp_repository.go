package user_view_aff_camp

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type userViewAffCampRepository struct {
	db *gorm.DB
}

func NewUserViewAffCampRepository(db *gorm.DB) interfaces.UserViewAffCampRepository {
	return &userViewAffCampRepository{
		db: db,
	}
}

func (r userViewAffCampRepository) CreateUserViewAffCamp(ctx context.Context, data *model.UserViewAffCamp) error {
	query := "INSERT INTO user_view_aff_camp (user_id, aff_camp_id) " +
		"VALUES (?, ?) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP"
	return r.db.Exec(query, data.UserId, data.AffCampId).Error
}

func (r userViewAffCampRepository) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) ([]model.UserViewAffCampComBrand, error) {
	var listUserViewAffCampComBrand []model.UserViewAffCampComBrand
	offset := (page - 1) * size
	err := r.db.Joins("AffCampComBrand").Where("user_id = ?", userId).
		Find(&listUserViewAffCampComBrand).Limit(size + 1).Offset(offset).Error
	return listUserViewAffCampComBrand, err
}

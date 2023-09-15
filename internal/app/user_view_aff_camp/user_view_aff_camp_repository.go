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

func (r userViewAffCampRepository) CreateUserViewAffCamp(ctx context.Context, data *model.UserViewAffCamp) error {
	query := "INSERT INTO user_view_brand (user_id, aff_camp_id) " +
		"VALUES (?, ?) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP"
	return r.db.Exec(query, data.UserId, data.AffCampId).Error
}

func NewUserViewAffCampRepository(db *gorm.DB) interfaces.UserViewAffCampRepository {
	return &userViewAffCampRepository{
		db: db,
	}
}

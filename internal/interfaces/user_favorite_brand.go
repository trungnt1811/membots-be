package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type UserFavoriteBrandRepository interface {
	GetListFavBrandByUserIdAndBrandIds(ctx context.Context, userId int64, brandIds []uint64) ([]model.UserFavoriteBrand, error)
}

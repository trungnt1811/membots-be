package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type UserFavoriteBrandRepository interface {
	GetListFavBrandByUserIdAndBrandIds(ctx context.Context, userId uint64, brandIds []uint) ([]model.UserFavoriteBrand, error)
}

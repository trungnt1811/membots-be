package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffBrandRepository interface {
	GetListFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error)
	UpdateCacheListFavouriteAffBrand(ctx context.Context) error
}

type AffBrandUCase interface {
	UpdateCacheListFavouriteAffBrand(ctx context.Context) error
}

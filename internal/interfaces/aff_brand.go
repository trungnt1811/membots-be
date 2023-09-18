package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffBrandRepository interface {
	GetListFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error)
}

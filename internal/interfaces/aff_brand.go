package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffBrandRepository interface {
	GetListCountFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error)
	UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error
}

type AffBrandUCase interface {
	UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error
	GetTopFavouriteAffBrand(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error)
}

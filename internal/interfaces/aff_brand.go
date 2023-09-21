package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffBrandRepository interface {
	GetListCountFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error)
	UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error
	GetListFavAffBrandByUserId(ctx context.Context, userId uint64, page, size int) ([]model.AffCampComFavBrand, error)
}

type AffBrandUCase interface {
	UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error
	GetTopFavouriteAffBrand(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error)
	GetListFavAffBrandByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error)
	GetMostCommissionAffCampaign(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error)
}

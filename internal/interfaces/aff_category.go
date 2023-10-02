package interfaces

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffCategoryUCase interface {
	GetAllCategory(ctx context.Context, page, size int) (dto.AffCategoryResponseDto, error)
	GetTopFavouriteAffBrand(ctx context.Context, categoryId uint, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error)
	GetMostCommissionAffCampaign(ctx context.Context, categoryId uint, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error)
}

type AffCategoryRepository interface {
	GetAllCategory(ctx context.Context, page, size int) ([]model.AffCategoryAndTotalCampaign, error)
}

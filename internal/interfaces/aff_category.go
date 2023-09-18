package interfaces

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffCategoryUCase interface {
	GetAllCategory(ctx context.Context, page, size int) (dto.AffCategoryResponseDto, error)
	GetAllAffCampaignInCategory(ctx context.Context, categoryId uint32, queryBy, order string, page, size int) (dto.AffCampaignAppDtoResponse, error)
}

type AffCategoryRepository interface {
	GetAllCategory(ctx context.Context, page, size int) ([]model.AffCategoryAndTotalCampaign, error)
	GetAllAffCampaignInCategory(ctx context.Context, categoryId uint32, orderBy string, page, size int) ([]model.AffCampaignLessApp, error)
}

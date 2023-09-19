package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffCampAppRepository interface {
	GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaignLessApp, error)
	GetAffCampaignById(ctx context.Context, id uint64) (model.AffCampaignApp, error)
	GetListAffCampaignByBrandIds(ctx context.Context, brandIds []uint64) ([]model.AffCampaignComBrand, error)
}

type AffCampAppUCase interface {
	GetAllAffCampaign(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error)
	GetAffCampaignById(ctx context.Context, id uint64, userId uint32) (dto.AffCampaignAppDto, error)
}

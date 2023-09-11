package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AppCampRepository interface {
	GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error)
	GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error)
}

type AppCampService interface {
	GetAllAffCampaign(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error)
	GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (dto.AffCampaignAppDto, error)
}

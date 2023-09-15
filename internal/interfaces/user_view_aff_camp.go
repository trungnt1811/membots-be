package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type UserViewAffCampRepository interface {
	CreateUserViewAffCamp(ctx context.Context, data *model.UserViewAffCamp) error
	GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) ([]model.UserViewAffCampComBrand, error)
}

type UserViewAffCampUCase interface {
	GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignComBrandDtoResponse, error)
}

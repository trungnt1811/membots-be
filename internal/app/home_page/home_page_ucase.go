package home_page

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type homePageUCase struct {
	AffBrandRepository          interfaces.AffBrandRepository
	AffCampAppRepository        interfaces.AffCampAppRepository
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
	UserViewAffCampRepository   interfaces.UserViewAffCampRepository
}

func NewHomePageUCase(affBrandRepository interfaces.AffBrandRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
	userViewAffCampRepository interfaces.UserViewAffCampRepository,
) interfaces.HomePageUCase {
	return &homePageUCase{
		AffBrandRepository:          affBrandRepository,
		AffCampAppRepository:        affCampAppRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
		UserViewAffCampRepository:   userViewAffCampRepository,
	}
}

func (s homePageUCase) GetHomePage(ctx context.Context, userId uint64) (dto.HomePageDto, error) {
	return dto.HomePageDto{}, nil
}

package aff_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affBrandUCase struct {
	AffBrandRepository   interfaces.AffBrandRepository
	AffCampAppRepository interfaces.AffCampAppRepository
}

func NewAffBrandUCase(
	affBrandRepository interfaces.AffBrandRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
) interfaces.AffBrandUCase {
	return &affBrandUCase{
		AffBrandRepository:   affBrandRepository,
		AffCampAppRepository: affCampAppRepository,
	}
}

func (s affBrandUCase) UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error {
	return s.AffBrandRepository.UpdateCacheListCountFavouriteAffBrand(ctx)
}

func (s affBrandUCase) GetTopFavouriteAffBrand(ctx context.Context, topFavourite int) ([]dto.AffCampaignComBrandDto, error) {
	listCountFavAffBrand, err := s.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
	if err != nil {
		return []dto.AffCampaignComBrandDto{}, err
	}

	// Get top favorited brands
	var brandIds []uint64
	for index, favCountAffBrand := range listCountFavAffBrand {
		if index >= topFavourite {
			break
		}
		brandIds = append(brandIds, favCountAffBrand.BrandId)
	}
	listFavAffBrand, err := s.AffCampAppRepository.GetListAffCampaignByBrandId(ctx, brandIds)
	if err != nil {
		return []dto.AffCampaignComBrandDto{}, err
	}

	var listAffCampaignComBrandDto []dto.AffCampaignComBrandDto
	for i := range listFavAffBrand {
		listAffCampaignComBrandDto = append(listAffCampaignComBrandDto, listFavAffBrand[i].ToAffCampaignComBrandDto())
	}
	return listAffCampaignComBrandDto, nil
}

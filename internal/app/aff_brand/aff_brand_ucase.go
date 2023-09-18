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

func (s affBrandUCase) GetTopFavouriteAffBrand(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listCountFavAffBrand, err := s.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	total := len(listCountFavAffBrand)

	nextPage := page
	if total > page*size {
		nextPage += 1
	}

	// Get top favorited brands
	var brandIds []uint64
	for index, favCountAffBrand := range listCountFavAffBrand {
		if index >= page*size || index < (page-1)*size {
			break
		}
		brandIds = append(brandIds, favCountAffBrand.BrandId)
	}
	listFavAffBrand, err := s.AffCampAppRepository.GetListAffCampaignByBrandIds(ctx, brandIds)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	var listAffCampaignComBrandDto []dto.AffCampaignLessDto
	for i := range listFavAffBrand {
		listAffCampaignComBrandDto = append(listAffCampaignComBrandDto, listFavAffBrand[i].ToAffCampaignLessDto())
	}
	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Total:    int64(total),
		Data:     listAffCampaignComBrandDto,
	}, nil
}

package aff_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affBrandUCase struct {
	AffBrandRepository          interfaces.AffBrandRepository
	AffCampAppRepository        interfaces.AffCampAppRepository
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
}

func NewAffBrandUCase(
	affBrandRepository interfaces.AffBrandRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
) interfaces.AffBrandUCase {
	return &affBrandUCase{
		AffBrandRepository:          affBrandRepository,
		AffCampAppRepository:        affCampAppRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
	}
}

func (s affBrandUCase) UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error {
	return s.AffBrandRepository.UpdateCacheListCountFavouriteAffBrand(ctx)
}

func (s affBrandUCase) GetTopFavouriteAffBrand(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listCountFavAffBrand, err := s.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	// Get top favorited brands
	var brandIds []uint
	for _, favCountAffBrand := range listCountFavAffBrand {
		brandIds = append(brandIds, favCountAffBrand.BrandId)
	}
	if len(brandIds) == 0 {
		return dto.AffCampaignAppDtoResponse{
			NextPage: page,
			Page:     page,
			Size:     size,
			Data:     nil,
		}, nil
	}
	listFavAffBrand, err := s.AffCampAppRepository.GetListAffCampaignByBrandIds(ctx, brandIds, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	// Check user's fav brand
	listUserFavBrand, err := s.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(ctx, userId, brandIds)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	favBrandCheck := make(map[uint]bool)
	for _, userFavBrand := range listUserFavBrand {
		favBrandCheck[userFavBrand.BrandId] = true
	}

	var listAffCampaignComBrandDto []dto.AffCampaignLessDto
	for i := range listFavAffBrand {
		if i >= size {
			break
		}
		listAffCampaignComBrandDto = append(listAffCampaignComBrandDto, listFavAffBrand[i].ToAffCampaignLessDto())
		listAffCampaignComBrandDto[i].Brand.IsFavorited = favBrandCheck[listAffCampaignComBrandDto[i].BrandId]
	}
	nextPage := page
	if len(listFavAffBrand) > size {
		nextPage += 1
	}
	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listAffCampaignComBrandDto,
	}, nil
}

func (s affBrandUCase) GetListFavAffBrandByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listFavAffBrand, err := s.AffBrandRepository.GetListFavAffBrandByUserId(ctx, userId, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	var listAffCampaignDto []dto.AffCampaignLessDto
	for i := range listFavAffBrand {
		listAffCampaignDto = append(listAffCampaignDto, listFavAffBrand[i].ToAffCampaignLessDto())
	}
	nextPage := page
	if len(listFavAffBrand) > size {
		nextPage += 1
	}
	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listAffCampaignDto,
	}, nil
}

func (s affBrandUCase) GetMostCommissionAffCampaign(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listAffCampaign, err := s.AffCampAppRepository.GetAllAffCampaign(ctx, interfaces.ListAffCampaignOrderByMostCommission, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	// Top favorited brands check
	listCountFavAffBrand, err := s.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	favTopBrandCheck := make(map[uint]bool)
	for _, countFavAffBrand := range listCountFavAffBrand {
		favTopBrandCheck[countFavAffBrand.BrandId] = true
	}

	// Use fav's brands check
	brandIds := make([]uint, 0)
	for _, affCampaign := range listAffCampaign {
		brandIds = append(brandIds, affCampaign.BrandId)
	}
	listUserFavBrand, err := s.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(
		ctx,
		userId,
		brandIds,
	)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	favBrandCheck := make(map[uint]bool)
	for _, userFavBrand := range listUserFavBrand {
		favBrandCheck[userFavBrand.BrandId] = true
	}

	// Parsing payload
	var listAffCampaignAppDto []dto.AffCampaignLessDto
	for i := range listAffCampaign {
		if i >= size {
			break
		}
		listAffCampaignAppDto = append(listAffCampaignAppDto, listAffCampaign[i].ToDto())
		listAffCampaignAppDto[i].Brand.IsFavorited = favBrandCheck[listAffCampaignAppDto[i].BrandId]
		listAffCampaignAppDto[i].Brand.IsTopFavorited = favTopBrandCheck[listAffCampaignAppDto[i].BrandId]
	}
	nextPage := page
	if len(listAffCampaign) > size {
		nextPage += 1
	}
	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listAffCampaignAppDto,
	}, nil
}

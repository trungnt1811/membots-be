package aff_camp_app

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type affCampAppUCase struct {
	AffCampAppRepository        interfaces.AffCampAppRepository
	AffBrandRepository          interfaces.AffBrandRepository
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
	Stream                      chan []*dto.UserViewAffCampDto
}

func NewAffCampAppUCase(
	affCampAppRepository interfaces.AffCampAppRepository,
	affBrandRepository interfaces.AffBrandRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
	stream chan []*dto.UserViewAffCampDto,
) interfaces.AffCampAppUCase {
	return &affCampAppUCase{
		AffCampAppRepository:        affCampAppRepository,
		AffBrandRepository:          affBrandRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
		Stream:                      stream,
	}
}

func (s affCampAppUCase) GetAffCampaignById(ctx context.Context, id uint64, userId uint32) (dto.AffCampaignAppDto, error) {
	affCampaign, err := s.AffCampAppRepository.GetAffCampaignById(ctx, id)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	if userId > 0 {
		payload := dto.UserViewAffCampDto{
			UserId:    userId,
			AffCampId: id,
		}
		uv := make([]*dto.UserViewAffCampDto, 0)
		uv = append(uv, &payload)
		s.Stream <- uv
	}

	listCountFavAffBrand, err := s.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	// Get top 10 fav brand
	favTopBrandCheck := make(map[uint64]bool)
	for index, countFavAffBrand := range listCountFavAffBrand {
		if index >= model.FavoritedBrandsInTop {
			break
		}
		favTopBrandCheck[countFavAffBrand.BrandId] = true
	}

	brandIds := make([]uint64, 0)
	brandIds = append(brandIds, affCampaign.BrandId)
	listUserFavBrand, err := s.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(
		ctx,
		uint64(userId),
		brandIds,
	)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	favBrandCheck := make(map[uint64]bool)
	for _, userFavBrand := range listUserFavBrand {
		favBrandCheck[userFavBrand.BrandId] = true
	}

	respone := affCampaign.ToAffCampaignAppDto()
	respone.Brand.IsTopFavorited = favTopBrandCheck[respone.BrandId]
	respone.Brand.IsFavorited = favBrandCheck[respone.BrandId]

	return respone, nil
}

func (s affCampAppUCase) GetAllAffCampaign(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listAffCampaign, err := s.AffCampAppRepository.GetAllAffCampaign(ctx, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	var listAffCampaignAppDto []dto.AffCampaignLessDto
	for i := range listAffCampaign {
		if i >= size {
			break
		}
		listAffCampaignAppDto = append(listAffCampaignAppDto, listAffCampaign[i].ToDto())
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

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
	ConvertPrice                interfaces.ConvertPriceHandler
}

func NewAffCampAppUCase(
	affCampAppRepository interfaces.AffCampAppRepository,
	affBrandRepository interfaces.AffBrandRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
	stream chan []*dto.UserViewAffCampDto,
	convertPrice interfaces.ConvertPriceHandler,
) interfaces.AffCampAppUCase {
	return &affCampAppUCase{
		AffCampAppRepository:        affCampAppRepository,
		AffBrandRepository:          affBrandRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
		Stream:                      stream,
		ConvertPrice:                convertPrice,
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
	favTopBrandCheck := make(map[uint]bool)
	for index, countFavAffBrand := range listCountFavAffBrand {
		if index >= model.FavoritedBrandsInTop {
			break
		}
		favTopBrandCheck[countFavAffBrand.BrandId] = true
	}

	brandIds := make([]uint, 0)
	brandIds = append(brandIds, affCampaign.BrandId)
	listUserFavBrand, err := s.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(
		ctx,
		uint64(userId),
		brandIds,
	)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	favBrandCheck := make(map[uint]bool)
	for _, userFavBrand := range listUserFavBrand {
		favBrandCheck[userFavBrand.BrandId] = true
	}

	response := affCampaign.ToAffCampaignAppDto()
	response.Brand.IsTopFavorited = favTopBrandCheck[response.BrandId]
	response.Brand.IsFavorited = favBrandCheck[response.BrandId]
	response.StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(ctx, affCampaign.Attributes)

	for i := range response.Attributes {
		tmp := model.AffCampaignAttribute{
			AttributeKey:   response.Attributes[i].AttributeKey,
			AttributeValue: response.Attributes[i].AttributeValue,
			AttributeType:  response.Attributes[i].AttributeType,
		}
		response.Attributes[i].AttributeValue = s.ConvertPrice.ConvertVndPriceToAstra(ctx, tmp)
		if tmp.AttributeType == "vnd" {
			response.Attributes[i].AttributeType = "ASA"
		}
	}

	return response, nil
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
		listAffCampaignAppDto[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(ctx, listAffCampaign[i].Attributes)
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

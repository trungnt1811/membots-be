package aff_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type affBrandUCase struct {
	AffBrandRepository          interfaces.AffBrandRepository
	AffCampAppRepository        interfaces.AffCampAppRepository
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
	ConvertPrice                interfaces.ConvertPriceHandler
}

func NewAffBrandUCase(
	affBrandRepository interfaces.AffBrandRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
	convertPrice interfaces.ConvertPriceHandler,
) interfaces.AffBrandUCase {
	return &affBrandUCase{
		AffBrandRepository:          affBrandRepository,
		AffCampAppRepository:        affCampAppRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
		ConvertPrice:                convertPrice,
	}
}

func (s affBrandUCase) UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error {
	return s.AffBrandRepository.UpdateCacheListCountFavouriteAffBrand(ctx)
}

func (s affBrandUCase) GetTopFavouriteAffBrand(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	// Get all attributes order by most commission
	listAffCampaignAttribute, err := s.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, interfaces.ListAffCampaignOrderByMostCommission)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	// Map only the most commision/aff campaign id
	campaignIdAtrributeMapping := make(map[uint64]model.AffCampaignAttribute)
	for _, attribute := range listAffCampaignAttribute {
		_, isExist := campaignIdAtrributeMapping[uint64(attribute.CampaignId)]
		if !isExist {
			campaignIdAtrributeMapping[uint64(attribute.CampaignId)] = attribute
		}
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
	for i, campaign := range listFavAffBrand {
		if i >= size {
			break
		}
		listAffCampaignComBrandDto = append(listAffCampaignComBrandDto, listFavAffBrand[i].ToAffCampaignLessDto())
		listAffCampaignComBrandDto[i].Brand.IsFavorited = favBrandCheck[listAffCampaignComBrandDto[i].BrandId]
		listAffCampaignComBrandDto[i].Brand.IsTopFavorited = favTopBrandCheck[listAffCampaignComBrandDto[i].BrandId]
		listAffCampaignComBrandDto[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.ID)]},
		)
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
	// Get all attributes order by most commission
	listAffCampaignAttribute, err := s.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, interfaces.ListAffCampaignOrderByMostCommission)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	// Map only the most commision/aff campaign id
	campaignIdAtrributeMapping := make(map[uint64]model.AffCampaignAttribute)
	for _, attribute := range listAffCampaignAttribute {
		_, isExist := campaignIdAtrributeMapping[uint64(attribute.CampaignId)]
		if !isExist {
			campaignIdAtrributeMapping[uint64(attribute.CampaignId)] = attribute
		}
	}

	listFavAffBrand, err := s.AffBrandRepository.GetListFavAffBrandByUserId(ctx, userId, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	total, err := s.AffBrandRepository.CountTotalFavAffBrandByUserId(ctx, userId)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	var listAffCampaignDto []dto.AffCampaignLessDto
	for i, campaign := range listFavAffBrand {
		listAffCampaignDto = append(listAffCampaignDto, listFavAffBrand[i].ToAffCampaignLessDto())
		listAffCampaignDto[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.ID)]},
		)
	}
	nextPage := page
	if len(listFavAffBrand) > size {
		nextPage += 1
	}
	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Total:    total,
		Data:     listAffCampaignDto,
	}, nil
}

func (s affBrandUCase) GetMostCommissionAffCampaign(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	// Get all attributes order by most commission
	listAffCampaignAttribute, err := s.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, interfaces.ListAffCampaignOrderByMostCommission)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	// Map only the most commision/aff campaign id
	campaignIdAtrributeMapping := make(map[uint64]model.AffCampaignAttribute)
	listAffCampaignId := make([]uint64, 0)
	for _, attribute := range listAffCampaignAttribute {
		_, isExist := campaignIdAtrributeMapping[uint64(attribute.CampaignId)]
		if !isExist {
			campaignIdAtrributeMapping[uint64(attribute.CampaignId)] = attribute
			listAffCampaignId = append(listAffCampaignId, uint64(attribute.CampaignId))
		}
	}

	// Get list aff campaign by ids
	listAffCampaign, err := s.AffCampAppRepository.GetListAffCampaignByIds(ctx, listAffCampaignId, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
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
		listAffCampaignAppDto = append(listAffCampaignAppDto, listAffCampaign[i].ToAffCampaignLessDto())
		listAffCampaignAppDto[i].Brand.IsFavorited = favBrandCheck[listAffCampaignAppDto[i].BrandId]
		listAffCampaignAppDto[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(listAffCampaignAppDto[i].ID)]},
		)
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

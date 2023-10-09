package category

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type categoryUCase struct {
	AffCategoryRepository       interfaces.AffCategoryRepository
	AffBrandRepository          interfaces.AffBrandRepository
	AffCampAppRepository        interfaces.AffCampAppRepository
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
	ConvertPrice                interfaces.ConvertPriceHandler
}

func (c *categoryUCase) GetMostCommissionAffCampaign(ctx context.Context, categoryId uint, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	// Get all attributes order by most commission
	listAffCampaignAttribute, err := c.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, interfaces.ListAffCampaignOrderByMostCommission)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	// Map only the most commission/aff campaign id
	campaignIdAtrributeMapping := make(map[uint64]model.AffCampaignAttribute)
	for _, attribute := range listAffCampaignAttribute {
		_, isExist := campaignIdAtrributeMapping[uint64(attribute.CampaignId)]
		if !isExist {
			campaignIdAtrributeMapping[uint64(attribute.CampaignId)] = attribute
		}
	}

	// Top favorited brands check
	listCountFavAffBrand, err := c.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
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
	listFavAffBrand, err := c.AffCampAppRepository.GetListAffCampaignByCategoryIdAndBrandIds(ctx, categoryId, brandIds, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	// Check user's fav brand
	listUserFavBrand, err := c.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(ctx, userId, brandIds)
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
		listAffCampaignComBrandDto[i].StellaMaxCom =
			c.ConvertPrice.GetStellaMaxCommission(ctx, []model.AffCampaignAttribute{campaignIdAtrributeMapping[listFavAffBrand[i].ID]})

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

func (c *categoryUCase) GetTopFavouriteAffBrand(ctx context.Context, categoryId uint, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	// Get all attributes order by most commission
	listAffCampaignAttribute, err := c.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, interfaces.ListAffCampaignOrderByMostCommission)
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

	// Get list aff campaign in category id
	listAffCampaign, err := c.AffCampAppRepository.GetAllAffCampaignInCategoryIdOrderByIds(ctx, categoryId, listAffCampaignId, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}

	// Top favorited brands check
	listCountFavAffBrand, err := c.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
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
	listUserFavBrand, err := c.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(
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
		listAffCampaignAppDto[i].StellaMaxCom = c.ConvertPrice.GetStellaMaxCommission(ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(listAffCampaign[i].ID)]})

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

func (c *categoryUCase) GetAllCategory(ctx context.Context, page, size int) (dto.AffCategoryResponseDto, error) {
	listCategory, err := c.AffCategoryRepository.GetAllCategory(ctx, page, size)
	if err != nil {
		return dto.AffCategoryResponseDto{}, err
	}
	var categoryDtos []dto.AffCategoryDto
	categoryIds := make([]uint64, size)
	for i := range listCategory {
		if i >= size {
			continue
		}
		categoryIds = append(categoryIds, listCategory[i].ID)
		categoryDtos = append(categoryDtos, listCategory[i].ToDto())
	}
	listCategoryWithAttribute, err := c.AffCategoryRepository.GetAttributeInCategories(ctx, categoryIds)
	mapCategoryAttribute := make(map[uint64][]model.AffCampaignAttribute)
	for _, attribute := range listCategoryWithAttribute {
		_, isExist := mapCategoryAttribute[attribute.ID]
		if !isExist {
			mapCategoryAttribute[attribute.ID] = make([]model.AffCampaignAttribute, 0)
		}
		mapCategoryAttribute[attribute.ID] = append(mapCategoryAttribute[attribute.ID], model.AffCampaignAttribute{
			AttributeValue: attribute.AttributeValue,
			AttributeType:  attribute.AttributeType,
		})
	}
	nextPage := page
	if len(listCategory) > size {
		nextPage = page + 1
	}
	for i := range categoryDtos {
		categoryDtos[i].StellaMaxCom = c.ConvertPrice.GetStellaMaxCommission(ctx, mapCategoryAttribute[categoryDtos[i].ID])
	}
	return dto.AffCategoryResponseDto{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     categoryDtos,
	}, nil
}

func NewAffCategoryUCase(repo interfaces.AffCategoryRepository,
	affBrandRepository interfaces.AffBrandRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
	convertPrice interfaces.ConvertPriceHandler,
) interfaces.AffCategoryUCase {
	return &categoryUCase{
		AffCategoryRepository:       repo,
		AffBrandRepository:          affBrandRepository,
		AffCampAppRepository:        affCampAppRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
		ConvertPrice:                convertPrice,
	}
}

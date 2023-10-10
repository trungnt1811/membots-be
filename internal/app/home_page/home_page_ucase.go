package home_page

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/chebyrash/promise"
)

type homePageUCase struct {
	AffBrandRepository          interfaces.AffBrandRepository
	AffCampAppRepository        interfaces.AffCampAppRepository
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
	UserViewAffCampRepository   interfaces.UserViewAffCampRepository
	ConvertPrice                interfaces.ConvertPriceHandler
}

func NewHomePageUCase(affBrandRepository interfaces.AffBrandRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
	userFavoriteBrandRepository interfaces.UserFavoriteBrandRepository,
	userViewAffCampRepository interfaces.UserViewAffCampRepository,
	convertPrice interfaces.ConvertPriceHandler,
) interfaces.HomePageUCase {
	return &homePageUCase{
		AffBrandRepository:          affBrandRepository,
		AffCampAppRepository:        affCampAppRepository,
		UserFavoriteBrandRepository: userFavoriteBrandRepository,
		UserViewAffCampRepository:   userViewAffCampRepository,
		ConvertPrice:                convertPrice,
	}
}

func (s homePageUCase) GetHomePage(ctx context.Context, userId uint64) (dto.HomePageDto, error) {
	// Get all attributes order by most commission
	listAffCampaignAttribute, err := s.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, interfaces.ListAffCampaignOrderByMostCommission)
	if err != nil {
		return dto.HomePageDto{}, err
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

	p1 := promise.New(func(resolve func([]model.UserViewAffCampComBrand), reject func(error)) {
		campaign, err := s.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, 1, 7-1)
		if err != nil {
			reject(err)
		} else {
			resolve(campaign)
		}
	})

	p2 := promise.New(func(resolve func([]model.AffCampComFavBrand), reject func(error)) {
		campaign, err := s.AffBrandRepository.GetListFavAffBrandByUserId(ctx, userId, 1, 7-1)
		if err != nil {
			reject(err)
		} else {
			resolve(campaign)
		}
	})

	p3 := promise.New(func(resolve func([]model.AffCampaignComBrand), reject func(error)) {
		listCountFavAffBrand, err := s.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
		if err != nil {
			reject(err)
		}

		// Get top favorited brands
		var brandIds []uint
		for _, favCountAffBrand := range listCountFavAffBrand {
			brandIds = append(brandIds, favCountAffBrand.BrandId)
		}
		if len(brandIds) == 0 {
			resolve(make([]model.AffCampaignComBrand, 0))
			return
		}

		campaign, err := s.AffCampAppRepository.GetListAffCampaignByBrandIds(ctx, brandIds, 1, 12-1)
		if err != nil {
			reject(err)
		} else {
			resolve(campaign)
		}
	})

	p4 := promise.New(func(resolve func([]model.AffCampaignComBrand), reject func(error)) {
		// Get list aff campaign by ids
		campaign, err := s.AffCampAppRepository.GetListAffCampaignByIds(ctx, listAffCampaignId, 1, 7-1)
		if err != nil {
			reject(err)
		} else {
			resolve(campaign)
		}
	})

	listCampaignP1, err := p1.Await(ctx)
	if err != nil {
		return dto.HomePageDto{}, err
	}

	listCampaignP2, err := p2.Await(ctx)
	if err != nil {
		return dto.HomePageDto{}, err
	}

	listCampaignP3, err := p3.Await(ctx)
	if err != nil {
		return dto.HomePageDto{}, err
	}

	listCampaignP4, err := p4.Await(ctx)
	if err != nil {
		return dto.HomePageDto{}, err
	}

	var listRecentlyVisited []dto.AffCampaignLessDto
	for i, campaign := range *listCampaignP1 {
		listRecentlyVisited = append(listRecentlyVisited, campaign.ToAffCampaignLessDto())
		listRecentlyVisited[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.AffCampId)]},
		)
	}

	var listFollowing []dto.AffCampaignLessDto
	for i, campaign := range *listCampaignP2 {
		listFollowing = append(listFollowing, campaign.ToAffCampaignLessDto())
		listFollowing[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.ID)]},
		)
	}

	var listTopFavorited []dto.AffCampaignLessDto
	for i, campaign := range *listCampaignP3 {
		listTopFavorited = append(listTopFavorited, campaign.ToAffCampaignLessDto())
		listTopFavorited[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.ID)]},
		)
	}

	var listMostCommission []dto.AffCampaignLessDto
	for i, campaign := range *listCampaignP4 {
		listMostCommission = append(listMostCommission, campaign.ToAffCampaignLessDto())
		listMostCommission[i].StellaMaxCom = s.ConvertPrice.GetStellaMaxCommission(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.ID)]},
		)
	}

	return dto.HomePageDto{
		RecentlyVisited: listRecentlyVisited,
		Following:       listFollowing,
		TopFavorited:    listTopFavorited,
		MostCommission:  listMostCommission,
	}, nil
}

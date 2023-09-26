package home_page

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	util "github.com/astraprotocol/affiliate-system/internal/util/commission"
	"github.com/chebyrash/promise"
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
	p1 := promise.New(func(resolve func([]model.UserViewAffCampComBrand), reject func(error)) {
		campaign, err := s.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, 1, 7)
		if err != nil {
			reject(err)
		} else {
			resolve(campaign)
		}
	})

	p2 := promise.New(func(resolve func([]model.AffCampComFavBrand), reject func(error)) {
		campaign, err := s.AffBrandRepository.GetListFavAffBrandByUserId(ctx, userId, 1, 7)
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

	var listRecentlyVisited []dto.AffCampaignLessDto
	for i, campaign := range *listCampaignP1 {
		listRecentlyVisited = append(listRecentlyVisited, campaign.ToAffCampaignLessDto())
		listRecentlyVisited[i].StellaMaxCom = util.GetStellaMaxCom(campaign.AffCampComBrand.Attributes)
	}

	var listFollowing []dto.AffCampaignLessDto
	for i, campaign := range *listCampaignP2 {
		listFollowing = append(listFollowing, campaign.ToAffCampaignLessDto())
		listFollowing[i].StellaMaxCom = util.GetStellaMaxCom(campaign.Attributes)
	}

	return dto.HomePageDto{
		RecentlyVisited: listRecentlyVisited,
		Following:       listFollowing,
	}, nil
}

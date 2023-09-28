package user_view_aff_camp

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type userViewAffCampUCase struct {
	UserViewAffCampRepository interfaces.UserViewAffCampRepository
	AffCampAppRepository      interfaces.AffCampAppRepository
	ConvertPrice              interfaces.ConvertPriceHandler
}

func NewUserViewAffCampUCase(userViewAffCampRepository interfaces.UserViewAffCampRepository,
	affCampAppRepository interfaces.AffCampAppRepository,
	convertPrice interfaces.ConvertPriceHandler) interfaces.UserViewAffCampUCase {
	return &userViewAffCampUCase{
		UserViewAffCampRepository: userViewAffCampRepository,
		AffCampAppRepository:      affCampAppRepository,
		ConvertPrice:              convertPrice,
	}
}

func (s *userViewAffCampUCase) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
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

	listUserViewAffCamp, err := s.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	var listAffCampComBrandDto []dto.AffCampaignLessDto
	for i, campaign := range listUserViewAffCamp {
		if i >= size {
			break
		}
		listAffCampComBrandDto = append(listAffCampComBrandDto, listUserViewAffCamp[i].ToAffCampaignLessDto())
		listAffCampComBrandDto[i].StellaMaxCom = s.ConvertPrice.ConvertVndPriceToAstra(
			ctx,
			[]model.AffCampaignAttribute{campaignIdAtrributeMapping[uint64(campaign.AffCampId)]},
		)
	}
	nextPage := page
	if len(listUserViewAffCamp) > size {
		nextPage += 1
	}
	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listAffCampComBrandDto,
	}, nil
}

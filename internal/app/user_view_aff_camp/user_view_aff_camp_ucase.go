package user_view_aff_camp

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type userViewAffCampUCase struct {
	UserViewAffCampRepository interfaces.UserViewAffCampRepository
	ConvertPrice              interfaces.ConvertPriceHandler
}

func NewUserViewAffCampUCase(repository interfaces.UserViewAffCampRepository,
	convertPrice interfaces.ConvertPriceHandler) interfaces.UserViewAffCampUCase {
	return &userViewAffCampUCase{
		UserViewAffCampRepository: repository,
		ConvertPrice:              convertPrice,
	}
}

func (s *userViewAffCampUCase) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listUserViewAffCamp, err := s.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	var listAffCampComBrandDto []dto.AffCampaignLessDto
	for i := range listUserViewAffCamp {
		if i >= size {
			break
		}
		listAffCampComBrandDto = append(listAffCampComBrandDto, listUserViewAffCamp[i].ToAffCampaignLessDto())
		listAffCampComBrandDto[i].StellaMaxCom = s.ConvertPrice.ConvertVndPriceToAstra(ctx, listUserViewAffCamp[i].AffCampComBrand.Attributes)
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

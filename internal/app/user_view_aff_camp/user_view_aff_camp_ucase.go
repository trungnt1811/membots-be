package user_view_aff_camp

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util "github.com/astraprotocol/affiliate-system/internal/util/commission"
)

type userViewAffCampUCase struct {
	UserViewAffCampRepository interfaces.UserViewAffCampRepository
}

func NewUserViewAffCampUCase(repository interfaces.UserViewAffCampRepository) interfaces.UserViewAffCampUCase {
	return &userViewAffCampUCase{
		UserViewAffCampRepository: repository,
	}
}

func (s userViewAffCampUCase) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignAppDtoResponse, error) {
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
		listAffCampComBrandDto[i].StellaMaxCom = util.GetStellaMaxCom(listUserViewAffCamp[i].AffCampComBrand.Attributes)
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

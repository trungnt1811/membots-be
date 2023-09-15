package user_view_aff_camp

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type userViewAffCampUCase struct {
	UserViewAffCampRepository interfaces.UserViewAffCampRepository
}

func NewUserViewAffCampUCase(repository interfaces.UserViewAffCampRepository) interfaces.UserViewAffCampUCase {
	return &userViewAffCampUCase{
		UserViewAffCampRepository: repository,
	}
}

func (s userViewAffCampUCase) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) (dto.AffCampaignComBrandDtoResponse, error) {
	listUserViewAffCamp, err := s.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, page, size)
	if err != nil {
		return dto.AffCampaignComBrandDtoResponse{}, err
	}
	var listAffCampComBrandDto []dto.AffCampaignComBrandDto
	for i := range listUserViewAffCamp {
		if i >= size {
			break
		}
		listAffCampComBrandDto = append(listAffCampComBrandDto, listUserViewAffCamp[i].ToAffCampaignComBrandDto())
	}
	nextPage := page
	if len(listUserViewAffCamp) > size {
		nextPage += 1
	}
	return dto.AffCampaignComBrandDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listAffCampComBrandDto,
	}, nil
}

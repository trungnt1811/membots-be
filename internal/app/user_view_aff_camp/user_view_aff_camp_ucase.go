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

func (s userViewAffCampUCase) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) (dto.UserViewAffCampComBrandDtoResponse, error) {
	listUserViewAffCamp, err := s.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, page, size)
	if err != nil {
		return dto.UserViewAffCampComBrandDtoResponse{}, err
	}
	var listUserViewAffCampDto []dto.UserViewAffCampComBrandDto
	for i := range listUserViewAffCamp {
		if i >= size {
			break
		}
		listUserViewAffCampDto = append(listUserViewAffCampDto, listUserViewAffCamp[i].ToUserViewAffCampComBrandDto())
	}
	nextPage := page
	if len(listUserViewAffCamp) > size {
		nextPage += 1
	}
	return dto.UserViewAffCampComBrandDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listUserViewAffCampDto,
	}, nil
}

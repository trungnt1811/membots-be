package user_view_aff_camp

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
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
		listAffCampComBrandDto[i].StellaMaxCom = s.getStellaMaxCom(listUserViewAffCamp[i].AffCampComBrand.Attributes)
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

func (s userViewAffCampUCase) getStellaMaxCom(attributes []model.AffCampaignAttribute) string {
	sort.Slice(attributes, func(i, j int) bool {
		if attributes[i].AttributeType != attributes[j].AttributeType {
			return model.AttributeTypePriorityMapping[attributes[i].AttributeType] < model.AttributeTypePriorityMapping[attributes[j].AttributeType]
		}
		switch strings.Compare(attributes[i].AttributeValue, attributes[j].AttributeValue) {
		case 1:
			return true
		default:
			return false
		}
	})
	if attributes[0].AttributeType == "percent" {
		return fmt.Sprint(attributes[0].AttributeValue, "%")
	} else {
		return fmt.Sprint(attributes[0].AttributeValue, " VND")
	}
}

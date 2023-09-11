package app_camp

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type appCampService struct {
	AppCampRepository interfaces.AppCampRepository
}

func (s appCampService) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (dto.AffCampaignAppDto, error) {
	affCampaign, err := s.AppCampRepository.GetAffCampaignByAccesstradeId(ctx, accesstradeId)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	affCampaignAppDto := affCampaign.ToAffCampaignAppDto()
	return affCampaignAppDto, nil
}

func (s appCampService) GetAllAffCampaign(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listAffCampaign, err := s.AppCampRepository.GetAllAffCampaign(ctx, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	var listAffCampaignAppDto []dto.AffCampaignAppDto
	for i := range listAffCampaign {
		if i >= size {
			break
		}
		listAffCampaignAppDto = append(listAffCampaignAppDto, listAffCampaign[i].ToAffCampaignAppDto())
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

func NewAppCampService(repository interfaces.AppCampRepository) interfaces.AppCampService {
	return &appCampService{
		AppCampRepository: repository,
	}
}

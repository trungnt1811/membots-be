package aff_camp_app

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affCampAppService struct {
	AffCampAppRepository interfaces.AffCampAppRepository
}

func NewAffCampAppService(repository interfaces.AffCampAppRepository) interfaces.AffCampAppService {
	return &affCampAppService{
		AffCampAppRepository: repository,
	}
}

func (s affCampAppService) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (dto.AffCampaignAppDto, error) {
	affCampaign, err := s.AffCampAppRepository.GetAffCampaignByAccesstradeId(ctx, accesstradeId)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	affCampaignAppDto := affCampaign.ToAffCampaignAppDto()
	return affCampaignAppDto, nil
}

func (s affCampAppService) GetAllAffCampaign(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listAffCampaign, err := s.AffCampAppRepository.GetAllAffCampaign(ctx, page, size)
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

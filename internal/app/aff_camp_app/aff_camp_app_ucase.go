package aff_camp_app

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affCampAppUCase struct {
	AffCampAppRepository interfaces.AffCampAppRepository
}

func NewAffCampAppUCase(repository interfaces.AffCampAppRepository) interfaces.AffCampAppUCase {
	return &affCampAppUCase{
		AffCampAppRepository: repository,
	}
}

func (s affCampAppUCase) GetAffCampaignById(ctx context.Context, id uint64) (dto.AffCampaignAppDto, error) {
	affCampaign, err := s.AffCampAppRepository.GetAffCampaignById(ctx, id)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	affCampaignAppDto := affCampaign.ToAffCampaignAppDto()
	return affCampaignAppDto, nil
}

func (s affCampAppUCase) GetAllAffCampaign(ctx context.Context, page, size int) (dto.AffCampaignAppDtoResponse, error) {
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

package aff_camp_app

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affCampAppUCase struct {
	AffCampAppRepository interfaces.AffCampAppRepository
	Stream               chan []*dto.UserViewAffCampDto
}

func NewAffCampAppUCase(repository interfaces.AffCampAppRepository, stream chan []*dto.UserViewAffCampDto) interfaces.AffCampAppUCase {
	return &affCampAppUCase{
		AffCampAppRepository: repository,
		Stream:               stream,
	}
}

func (s affCampAppUCase) GetAffCampaignById(ctx context.Context, id uint64, userId uint32) (dto.AffCampaignAppDto, error) {
	affCampaign, err := s.AffCampAppRepository.GetAffCampaignById(ctx, id)
	if err != nil {
		return dto.AffCampaignAppDto{}, err
	}
	if userId > 0 {
		payload := dto.UserViewAffCampDto{
			UserId:    userId,
			AffCampId: id,
		}
		uv := make([]*dto.UserViewAffCampDto, 0)
		uv = append(uv, &payload)
		s.Stream <- uv
	}
	return affCampaign.ToAffCampaignAppDto(), nil
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

package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type campaignUCase struct {
	Repo interfaces.ConsoleCampRepository
}

func (c *campaignUCase) GetAllCampaign(status []string, page, size int) (dto.AffCampaignDtoResponse, error) {
	listAffCampaign, err := c.Repo.GetAllCampaign(status, page, size)
	if err != nil {
		return dto.AffCampaignDtoResponse{}, err
	}
	totalCampaign, err := c.Repo.CountCampaign(status)
	if err != nil {
		return dto.AffCampaignDtoResponse{}, err
	}
	var listAffCampaignDto []dto.AffCampaignDto
	for i, campaign := range listAffCampaign {
		if i >= size {
			continue
		}
		listAffCampaignDto = append(listAffCampaignDto, campaign.ToDto())
	}
	nextPage := page
	if len(listAffCampaign) > size {
		nextPage += 1
	}
	return dto.AffCampaignDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Total:    totalCampaign,
		Data:     listAffCampaignDto,
	}, nil
}

func NewCampaignUCase(repo interfaces.ConsoleCampRepository) interfaces.ConsoleCampUCase {
	return &campaignUCase{
		Repo: repo,
	}
}

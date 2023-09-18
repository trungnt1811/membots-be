package campaign

import (
	"encoding/json"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"gorm.io/datatypes"
	"strings"
)

type campaignUCase struct {
	Repo interfaces.ConsoleCampRepository
}

func (c *campaignUCase) GetCampaignById(id uint) (dto.AffCampaignDto, error) {
	affCampaign, err := c.Repo.GetCampaignById(id)
	if err != nil {
		return dto.AffCampaignDto{}, err
	}
	return affCampaign.ToDto(), nil
}

func (c *campaignUCase) UpdateCampaign(id uint, campaign dto.AffCampaignAppDto) error {
	updates := make(map[string]interface{})
	if campaign.StellaDescription != nil {
		b, err := json.Marshal(campaign.StellaDescription)
		if err != nil {
			return err
		}
		updates["stella_description"] = datatypes.JSON(b)
	}
	if len(strings.TrimSpace(campaign.Name)) > 0 {
		updates["name"] = campaign.Name
	}
	if len(strings.TrimSpace(campaign.Thumbnail)) > 0 {
		updates["thumbnail"] = campaign.Thumbnail
	}
	if len(strings.TrimSpace(campaign.Url)) > 0 {
		updates["url"] = campaign.Url
	}
	if len(strings.TrimSpace(campaign.StellaMaxCom)) > 0 {
		updates["stella_max_com"] = campaign.StellaMaxCom
	}
	if campaign.StartTime != nil {
		updates["start_time"] = campaign.StartTime
	}
	if campaign.EndTime != nil {
		updates["end_time"] = campaign.EndTime
	}
	if campaign.CategoryId > 0 {
		updates["category_id"] = campaign.CategoryId
	}
	if campaign.BrandId > 0 {
		updates["brand_id"] = campaign.BrandId
	}
	if len(strings.TrimSpace(campaign.StellaStatus)) > 0 {
		updates["stella_status"] = campaign.StellaStatus
	}
	return c.Repo.UpdateCampaign(id, updates)
}

func (c *campaignUCase) GetAllCampaign(status []string, q string, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	listAffCampaign, err := c.Repo.GetAllCampaign(status, q, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	totalCampaign, err := c.Repo.CountCampaign(status, q)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	var listAffCampaignDto []dto.AffCampaignLessDto
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
	return dto.AffCampaignAppDtoResponse{
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

package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	model2 "github.com/astraprotocol/affiliate-system/internal/model"
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
		updates["stella_description"] = campaign.StellaDescription
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

	if len(campaign.Attributes) > 0 {
		var attributes []model2.AffCampaignAttribute
		for _, attribute := range campaign.Attributes {
			attributes = append(attributes, model2.AffCampaignAttribute{
				ID:             attribute.ID,
				CampaignId:     id,
				AttributeKey:   attribute.AttributeKey,
				AttributeValue: attribute.AttributeValue,
				AttributeType:  attribute.AttributeType,
			})
		}
		err := c.Repo.UpdateCampaignAttribute(id, attributes)
		if err != nil {
			return err
		}
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

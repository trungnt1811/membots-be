package campaign

import (
	"errors"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"strings"
)

type bannerUCase struct {
	Repo     interfaces.ConsoleBannerRepository
	CampRepo interfaces.ConsoleCampRepository
}

func (c *bannerUCase) CreateBanner(banner *dto.AffBannerDto) (dto.AffBannerDto, error) {
	_, err := c.CampRepo.GetCampaignLessByAccessTradeId(banner.AccessTradeId)
	if err != nil {
		return dto.AffBannerDto{}, errors.New("accesstrade_id not found")
	}
	bannerData, err := c.Repo.CreateBanner(&model.AffBanner{
		Name:          banner.Name,
		Thumbnail:     banner.Thumbnail,
		Url:           banner.Url,
		AccessTradeId: banner.AccessTradeId,
		Status:        "inactive",
	})
	return bannerData.ToDto(), nil
}

func (c *bannerUCase) GetBannerById(id uint) (dto.AffBannerDto, error) {
	affCampaign, err := c.Repo.GetBannerById(id)
	if err != nil {
		return dto.AffBannerDto{}, err
	}
	return affCampaign.ToDto(), nil
}

func (c *bannerUCase) UpdateBanner(id uint, campaign *dto.AffBannerDto) error {
	updates := make(map[string]interface{})
	if len(strings.TrimSpace(campaign.Name)) > 0 {
		updates["name"] = campaign.Name
	}
	if len(strings.TrimSpace(campaign.Thumbnail)) > 0 {
		updates["thumbnail"] = campaign.Thumbnail
	}
	if len(strings.TrimSpace(campaign.Url)) > 0 {
		updates["url"] = campaign.Url
	}
	if len(strings.TrimSpace(campaign.Status)) > 0 {
		updates["status"] = campaign.Status
	}
	return c.Repo.UpdateBanner(id, updates)
}

func (c *bannerUCase) GetAllBanner(status []string, page, size int) (dto.AffBannerDtoResponse, error) {
	listAffBanner, err := c.Repo.GetAllBanner(status, page, size)
	if err != nil {
		return dto.AffBannerDtoResponse{}, err
	}
	totalCampaign, err := c.Repo.CountBanner(status)
	if err != nil {
		return dto.AffBannerDtoResponse{}, err
	}
	var listAffBannerDto []dto.AffBannerDto
	for i, banner := range listAffBanner {
		if i >= size {
			continue
		}
		listAffBannerDto = append(listAffBannerDto, banner.ToDto())
	}
	nextPage := page
	if len(listAffBanner) > size {
		nextPage += 1
	}
	return dto.AffBannerDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Total:    totalCampaign,
		Data:     listAffBannerDto,
	}, nil
}

func NewBannerUCase(repo interfaces.ConsoleBannerRepository, campRepo interfaces.ConsoleCampRepository) interfaces.ConsoleBannerUCase {
	return &bannerUCase{
		Repo:     repo,
		CampRepo: campRepo,
	}
}

package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type bannerUCase struct {
	Repo interfaces.AppBannerRepository
}

func (c *bannerUCase) GetBannerById(id uint) (dto.AffBannerDto, error) {
	affCampaign, err := c.Repo.GetBannerById(id)
	if err != nil {
		return dto.AffBannerDto{}, err
	}
	return affCampaign.ToDto(), nil
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

func NewBannerUCase(repo interfaces.AppBannerRepository) interfaces.AppBannerUCase {
	return &bannerUCase{
		Repo: repo,
	}
}

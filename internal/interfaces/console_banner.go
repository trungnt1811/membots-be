package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type ConsoleBannerRepository interface {
	GetAllBanner(status []string, page, size int) ([]model.AffBanner, error)
	GetBannerById(id uint) (model.AffBanner, error)
	CountBanner(status []string) (int64, error)
	UpdateBanner(id uint, updates map[string]interface{}) error
	CreateBanner(banner *model.AffBanner) (model.AffBanner, error)
}

type ConsoleBannerUCase interface {
	GetAllBanner(status []string, page, size int) (dto.AffBannerDtoResponse, error)
	UpdateBanner(id uint, campaign *dto.AffBannerDto) error
	GetBannerById(id uint) (dto.AffBannerDto, error)
	CreateBanner(banner *dto.AffBannerDto) (dto.AffBannerDto, error)
}

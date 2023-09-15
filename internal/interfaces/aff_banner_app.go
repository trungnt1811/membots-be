package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AppBannerRepository interface {
	GetAllBanner(status []string, page, size int) ([]model.AffBanner, error)
	GetBannerById(id uint) (model.AffBanner, error)
	CountBanner(status []string) (int64, error)
}

type AppBannerUCase interface {
	GetAllBanner(status []string, page, size int) (dto.AffBannerDtoResponse, error)
	GetBannerById(id uint) (dto.AffBannerDto, error)
}

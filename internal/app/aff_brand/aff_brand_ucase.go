package aff_brand

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affBrandUCase struct {
	AffBrandRepository interfaces.AffBrandRepository
}

func NewAffBrandUCase(repository interfaces.AffBrandRepository) interfaces.AffBrandUCase {
	return &affBrandUCase{
		AffBrandRepository: repository,
	}
}

func (s affBrandUCase) UpdateCacheListFavouriteAffBrand(ctx context.Context) error {
	return s.AffBrandRepository.UpdateCacheListFavouriteAffBrand(ctx)
}

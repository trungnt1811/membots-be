package aff_brand

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffBrand = "aff_brand_"
const cacheTimeAffBrand = 30 * time.Minute

type affBrandCache struct {
	AffBrandRepository interfaces.AffBrandRepository
	Cache              caching.Repository
}

func NewAffBrandCacheRepository(repo interfaces.AffBrandRepository,
	cache caching.Repository,
) interfaces.AffBrandRepository {
	return &affBrandCache{
		AffBrandRepository: repo,
		Cache:              cache,
	}
}

func (c affBrandCache) GetListFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error) {
	key := &caching.Keyer{Raw: keyPrefixAffBrand + "GetListFavouriteAffBrand"}
	var listFavouriteAffBrand []model.TotalFavoriteBrand
	err := c.Cache.RetrieveItem(key, &listFavouriteAffBrand)
	if err != nil {
		// cache miss
		listFavouriteAffBrand, err = c.AffBrandRepository.GetListFavouriteAffBrand(ctx)
		if err != nil {
			return listFavouriteAffBrand, err
		}
		if err = c.Cache.SaveItem(key, listFavouriteAffBrand, cacheTimeAffBrand); err != nil {
			return listFavouriteAffBrand, err
		}
	}
	return listFavouriteAffBrand, nil
}

func (c affBrandCache) UpdateCacheListFavouriteAffBrand(ctx context.Context) error {
	key := &caching.Keyer{Raw: keyPrefixAffBrand + "GetListFavouriteAffBrand"}
	listFavouriteAffBrand, err := c.AffBrandRepository.GetListFavouriteAffBrand(ctx)
	if err != nil {
		return err
	}
	if err = c.Cache.SaveItem(key, listFavouriteAffBrand, cacheTimeAffBrand); err != nil {
		return err
	}
	return nil
}

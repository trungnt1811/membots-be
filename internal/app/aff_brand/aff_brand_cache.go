package aff_brand

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffBrand = "aff_brand_"
const cacheTimeLongAffBrand = 30 * time.Minute
const cacheTimeAffBrand = 3 * time.Second

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

func (c affBrandCache) GetListCountFavouriteAffBrand(ctx context.Context) ([]model.TotalFavoriteBrand, error) {
	key := &caching.Keyer{Raw: keyPrefixAffBrand + "GetListCountFavouriteAffBrand"}
	var listFavouriteAffBrand []model.TotalFavoriteBrand
	err := c.Cache.RetrieveItem(key, &listFavouriteAffBrand)
	if err != nil {
		// cache miss
		listFavouriteAffBrand, err = c.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
		if err != nil {
			return listFavouriteAffBrand, err
		}
		if err = c.Cache.SaveItem(key, listFavouriteAffBrand, cacheTimeLongAffBrand); err != nil {
			return listFavouriteAffBrand, err
		}
	}
	return listFavouriteAffBrand, nil
}

func (c affBrandCache) UpdateCacheListCountFavouriteAffBrand(ctx context.Context) error {
	key := &caching.Keyer{Raw: keyPrefixAffBrand + "GetListCountFavouriteAffBrand"}
	listFavouriteAffBrand, err := c.AffBrandRepository.GetListCountFavouriteAffBrand(ctx)
	if err != nil {
		return err
	}
	if err = c.Cache.SaveItem(key, listFavouriteAffBrand, cacheTimeLongAffBrand); err != nil {
		return err
	}
	return nil
}

func (c affBrandCache) GetListFavAffBrandByUserId(ctx context.Context, userId uint64, page, size int) ([]model.AffCampComFavBrand, error) {
	key := &caching.Keyer{Raw: keyPrefixAffBrand + fmt.Sprint("GetListFavAffBrandByUserId_", userId, "_", page, "_", size)}
	var listAffCampComFavBrand []model.AffCampComFavBrand
	err := c.Cache.RetrieveItem(key, &listAffCampComFavBrand)
	if err != nil {
		// cache miss
		listAffCampComFavBrand, err = c.AffBrandRepository.GetListFavAffBrandByUserId(ctx, userId, page, size)
		if err != nil {
			return listAffCampComFavBrand, err
		}
		if err = c.Cache.SaveItem(key, listAffCampComFavBrand, cacheTimeAffBrand); err != nil {
			return listAffCampComFavBrand, err
		}
	}
	return listAffCampComFavBrand, nil
}

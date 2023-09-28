package aff_camp_app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffCampApp = "aff_camp_app_"
const cacheTimeAffCampApp = 3 * time.Second

type affCampAppCache struct {
	AffCampAppRepository interfaces.AffCampAppRepository
	Cache                caching.Repository
}

func NewAffCampAppCacheRepository(repo interfaces.AffCampAppRepository,
	cache caching.Repository,
) interfaces.AffCampAppRepository {
	return &affCampAppCache{
		AffCampAppRepository: repo,
		Cache:                cache,
	}
}

func (c affCampAppCache) GetAllAffCampaignInCategoryId(ctx context.Context, categoryId uint, orderBy string, page, size int) ([]model.AffCampaignLessApp, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAllAffCampaignInCategoryId_", categoryId, "_", orderBy, "_", page, "_", size)}
	var listAffCampaign []model.AffCampaignLessApp
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetAllAffCampaignInCategoryId(ctx, categoryId, orderBy, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampApp); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c affCampAppCache) GetListAffCampaignByCategoryIdAndBrandIds(ctx context.Context, categoryId uint, brandIds []uint, page, size int) ([]model.AffCampaignComBrand, error) {
	s, _ := json.Marshal(brandIds)
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetListAffCampaignByCategoryIdAndBrandIds_", categoryId, "_", string(s), "_", page, "_", size)}
	var listAffCampaign []model.AffCampaignComBrand
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetListAffCampaignByCategoryIdAndBrandIds(ctx, categoryId, brandIds, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampApp); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c affCampAppCache) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaignLessApp, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAllAffCampaign_", page, "_", size)}
	var listAffCampaign []model.AffCampaignLessApp
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetAllAffCampaign(ctx, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampApp); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c affCampAppCache) GetAffCampaignById(ctx context.Context, id uint64) (model.AffCampaignApp, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAffCampaignById_", id)}
	var affCampaign model.AffCampaignApp
	err := c.Cache.RetrieveItem(key, &affCampaign)
	if err != nil {
		// cache miss
		affCampaign, err = c.AffCampAppRepository.GetAffCampaignById(ctx, id)
		if err != nil {
			return affCampaign, err
		}
		if err = c.Cache.SaveItem(key, affCampaign, cacheTimeAffCampApp); err != nil {
			return affCampaign, err
		}
	}
	return affCampaign, nil
}

func (c affCampAppCache) GetListAffCampaignByBrandIds(ctx context.Context, brandIds []uint, page, size int) ([]model.AffCampaignComBrand, error) {
	s, _ := json.Marshal(brandIds)
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAffCampaignById_", string(s), "_", page, "_", size)}
	var listAffCampaign []model.AffCampaignComBrand
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetListAffCampaignByBrandIds(ctx, brandIds, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampApp); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c affCampAppCache) GetAllAffCampaignAttribute(ctx context.Context, orderBy string) ([]model.AffCampaignAttribute, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAllAffCampaignAttribute_", orderBy)}
	var listAffCampaignAttribute []model.AffCampaignAttribute
	err := c.Cache.RetrieveItem(key, &listAffCampaignAttribute)
	if err != nil {
		// cache miss
		listAffCampaignAttribute, err = c.AffCampAppRepository.GetAllAffCampaignAttribute(ctx, orderBy)
		if err != nil {
			return listAffCampaignAttribute, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaignAttribute, cacheTimeAffCampApp); err != nil {
			return listAffCampaignAttribute, err
		}
	}
	return listAffCampaignAttribute, nil
}

func (c affCampAppCache) GetListAffCampaignByIds(ctx context.Context, ids []uint64, page, size int) ([]model.AffCampaignComBrand, error) {
	s, _ := json.Marshal(ids)
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetListAffCampaignByIds_", string(s), "_", page, "_", size)}
	var listAffCampaign []model.AffCampaignComBrand
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetListAffCampaignByIds(ctx, ids, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampApp); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

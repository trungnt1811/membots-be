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

func (c affCampAppCache) GetAllAffCampaign(ctx context.Context, orderBy string, page, size int) ([]model.AffCampaignLessApp, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAllAffCampaign_", orderBy, "_", page, "_", size)}
	var listAffCampaign []model.AffCampaignLessApp
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetAllAffCampaign(ctx, orderBy, page, size)
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

func (c affCampAppCache) GetListAffCampaignByBrandIds(ctx context.Context, brandIds []uint64) ([]model.AffCampaignComBrand, error) {
	s, _ := json.Marshal(brandIds)
	key := &caching.Keyer{Raw: keyPrefixAffCampApp + fmt.Sprint("GetAffCampaignById_", string(s))}
	var listAffCampaign []model.AffCampaignComBrand
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetListAffCampaignByBrandIds(ctx, brandIds)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampApp); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

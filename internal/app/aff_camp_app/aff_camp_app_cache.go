package aff_camp_app

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffCampaign = "aff_camp_app_"
const cacheTimeAffCampaign = 3 * time.Second

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

func (c affCampAppCache) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaignApp, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAllAffCampaign_", page, "_", size)}
	var listAffCampaign []model.AffCampaignApp
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampAppRepository.GetAllAffCampaign(ctx, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampaign); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c affCampAppCache) GetAffCampaignById(ctx context.Context, accesstradeId uint64) (model.AffCampaignApp, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAffCampaignById_", accesstradeId)}
	var affCampaign model.AffCampaignApp
	err := c.Cache.RetrieveItem(key, &affCampaign)
	if err != nil {
		// cache miss
		affCampaign, err = c.AffCampAppRepository.GetAffCampaignById(ctx, accesstradeId)
		if err != nil {
			return affCampaign, err
		}
		if err = c.Cache.SaveItem(key, affCampaign, cacheTimeAffCampaign); err != nil {
			return affCampaign, err
		}
	}
	return affCampaign, nil
}

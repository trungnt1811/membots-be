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

func (c affCampAppCache) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAllAffCampaign_", page, "_", size)}
	var listAffCampaign []model.AffCampaign
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

func (c affCampAppCache) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAffCampaignByAccesstradeId_", accesstradeId)}
	var affCampaign model.AffCampaign
	err := c.Cache.RetrieveItem(key, &affCampaign)
	if err != nil {
		// cache miss
		affCampaign, err = c.AffCampAppRepository.GetAffCampaignByAccesstradeId(ctx, accesstradeId)
		if err != nil {
			return affCampaign, err
		}
		if err = c.Cache.SaveItem(key, affCampaign, cacheTimeAffCampaign); err != nil {
			return affCampaign, err
		}
	}
	return affCampaign, nil
}

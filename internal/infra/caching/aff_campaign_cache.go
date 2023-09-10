package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffCampaign = "aff_campaign_"
const cacheTimeAffCampaign = 3 * time.Second

type affCampaignCache struct {
	AffCampaignRepository interfaces.AffCampaignRepository
	Cache                 cachingRepository
}

func (c affCampaignCache) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error) {
	key := &Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAllAffCampaign_", page, "_", size)}
	var listAffCampaign []model.AffCampaign
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AffCampaignRepository.GetAllAffCampaign(ctx, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampaign); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c affCampaignCache) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error) {
	key := &Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAffCampaignByAccesstradeId_", accesstradeId)}
	var affCampaign model.AffCampaign
	err := c.Cache.RetrieveItem(key, &affCampaign)
	if err != nil {
		// cache miss
		affCampaign, err = c.AffCampaignRepository.GetAffCampaignByAccesstradeId(ctx, accesstradeId)
		if err != nil {
			return affCampaign, err
		}
		if err = c.Cache.SaveItem(key, affCampaign, cacheTimeAffCampaign); err != nil {
			return affCampaign, err
		}
	}
	return affCampaign, nil
}

func NewAffCampaignCacheRepository(repo interfaces.AffCampaignRepository,
	cache cachingRepository,
) interfaces.AffCampaignRepository {
	return &affCampaignCache{
		AffCampaignRepository: repo,
		Cache:                 cache,
	}
}

package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffCampaign = "app_camp_"
const cacheTimeAffCampaign = 3 * time.Second

type appCampCache struct {
	AppCampRepository interfaces.AppCampRepository
	Cache             cachingRepository
}

func (c appCampCache) GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error) {
	key := &Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAllAffCampaign_", page, "_", size)}
	var listAffCampaign []model.AffCampaign
	err := c.Cache.RetrieveItem(key, &listAffCampaign)
	if err != nil {
		// cache miss
		listAffCampaign, err = c.AppCampRepository.GetAllAffCampaign(ctx, page, size)
		if err != nil {
			return listAffCampaign, err
		}
		if err = c.Cache.SaveItem(key, listAffCampaign, cacheTimeAffCampaign); err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (c appCampCache) GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error) {
	key := &Keyer{Raw: keyPrefixAffCampaign + fmt.Sprint("GetAffCampaignByAccesstradeId_", accesstradeId)}
	var affCampaign model.AffCampaign
	err := c.Cache.RetrieveItem(key, &affCampaign)
	if err != nil {
		// cache miss
		affCampaign, err = c.AppCampRepository.GetAffCampaignByAccesstradeId(ctx, accesstradeId)
		if err != nil {
			return affCampaign, err
		}
		if err = c.Cache.SaveItem(key, affCampaign, cacheTimeAffCampaign); err != nil {
			return affCampaign, err
		}
	}
	return affCampaign, nil
}

func NewAppCampCacheRepository(repo interfaces.AppCampRepository,
	cache cachingRepository,
) interfaces.AppCampRepository {
	return &appCampCache{
		AppCampRepository: repo,
		Cache:             cache,
	}
}

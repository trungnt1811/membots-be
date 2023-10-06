package category

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixAffCategoryApp = "aff_category_app_"
const cacheTimeAffCategoryApp = 30 * time.Second

type affCategoryAppCache struct {
	AffCategoryRepository interfaces.AffCategoryRepository
	Cache                 caching.Repository
}

func (c affCategoryAppCache) GetAllCategory(ctx context.Context, page, size int) ([]model.AffCategoryAndTotalCampaign, error) {
	key := &caching.Keyer{Raw: keyPrefixAffCategoryApp + fmt.Sprint("GetAllCategory_", page, "_", size)}
	var listCategory []model.AffCategoryAndTotalCampaign
	err := c.Cache.RetrieveItem(key, &listCategory)
	if err != nil {
		// cache miss
		listCategory, err = c.AffCategoryRepository.GetAllCategory(ctx, page, size)
		if err != nil {
			return listCategory, err
		}
		if err = c.Cache.SaveItem(key, listCategory, cacheTimeAffCategoryApp); err != nil {
			return listCategory, err
		}
	}
	return listCategory, nil
}

func (c affCategoryAppCache) GetAttributeInCategories(ctx context.Context, ids []uint64) ([]model.CategoryWithCommissionAttribute, error) {
	s, _ := json.Marshal(ids)
	key := &caching.Keyer{Raw: keyPrefixAffCategoryApp + fmt.Sprint("GetAttributeInCategories_", string(s))}
	var listCategory []model.CategoryWithCommissionAttribute
	err := c.Cache.RetrieveItem(key, &listCategory)
	if err != nil {
		// cache miss
		listCategory, err = c.AffCategoryRepository.GetAttributeInCategories(ctx, ids)
		if err != nil {
			return listCategory, err
		}
		if err = c.Cache.SaveItem(key, listCategory, cacheTimeAffCategoryApp); err != nil {
			return listCategory, err
		}
	}
	return listCategory, nil
}

func NewAffCategoryAppCacheRepository(repo interfaces.AffCategoryRepository,
	cache caching.Repository,
) interfaces.AffCategoryRepository {
	return &affCategoryAppCache{
		AffCategoryRepository: repo,
		Cache:                 cache,
	}
}

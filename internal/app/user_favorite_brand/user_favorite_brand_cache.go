package user_favorite_brand

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixUserFavoriteBrand = "user_favorite_brand_"
const cacheTimeUserFavoriteBrand = 3 * time.Second

type userFavoriteBrandCache struct {
	UserFavoriteBrandRepository interfaces.UserFavoriteBrandRepository
	Cache                       caching.Repository
}

func NewUserFavoriteBrandCacheRepository(repo interfaces.UserFavoriteBrandRepository,
	cache caching.Repository,
) interfaces.UserFavoriteBrandRepository {
	return &userFavoriteBrandCache{
		UserFavoriteBrandRepository: repo,
		Cache:                       cache,
	}
}

func (c userFavoriteBrandCache) GetListFavBrandByUserIdAndBrandIds(ctx context.Context, userId uint64, brandIds []uint64) ([]model.UserFavoriteBrand, error) {
	s, _ := json.Marshal(brandIds)
	key := &caching.Keyer{Raw: keyPrefixUserFavoriteBrand + fmt.Sprint("GetListFavBrandByUserIdAndBrandIds_", userId, "_", string(s))}
	var listUserFavoriteBrand []model.UserFavoriteBrand
	err := c.Cache.RetrieveItem(key, &listUserFavoriteBrand)
	if err != nil {
		// cache miss
		listUserFavoriteBrand, err = c.UserFavoriteBrandRepository.GetListFavBrandByUserIdAndBrandIds(ctx, userId, brandIds)
		if err != nil {
			return listUserFavoriteBrand, err
		}
		if err = c.Cache.SaveItem(key, listUserFavoriteBrand, cacheTimeUserFavoriteBrand); err != nil {
			return listUserFavoriteBrand, err
		}
	}
	return listUserFavoriteBrand, nil
}

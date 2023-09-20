package user_view_aff_camp

import (
	"context"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

const keyPrefixUserViewAffCamp = "user_view_aff_camp_"
const cacheTimeUserViewAffCamp = 3 * time.Second

type userViewAffCampCache struct {
	UserViewAffCampRepository interfaces.UserViewAffCampRepository
	Cache                     caching.Repository
}

func NewUserViewAffCampCacheRepository(repo interfaces.UserViewAffCampRepository,
	cache caching.Repository,
) interfaces.UserViewAffCampRepository {
	return &userViewAffCampCache{
		UserViewAffCampRepository: repo,
		Cache:                     cache,
	}
}

func (c userViewAffCampCache) CreateUserViewAffCamp(ctx context.Context, data *model.UserViewAffCamp) error {
	return c.UserViewAffCampRepository.CreateUserViewAffCamp(ctx, data)
}

func (c userViewAffCampCache) GetListUserViewAffCampByUserId(ctx context.Context, userId uint64, page, size int) ([]model.UserViewAffCampComBrand, error) {
	key := &caching.Keyer{Raw: keyPrefixUserViewAffCamp + fmt.Sprint("GetListUserViewAffCampByUserId_", userId, "_", page, "_", size)}
	var listUserViewAffCamp []model.UserViewAffCampComBrand
	err := c.Cache.RetrieveItem(key, &listUserViewAffCamp)
	if err != nil {
		// cache miss
		listUserViewAffCamp, err = c.UserViewAffCampRepository.GetListUserViewAffCampByUserId(ctx, userId, page, size)
		if err != nil {
			return listUserViewAffCamp, err
		}
		if err = c.Cache.SaveItem(key, listUserViewAffCamp, cacheTimeUserViewAffCamp); err != nil {
			return listUserViewAffCamp, err
		}
	}
	return listUserViewAffCamp, nil
}

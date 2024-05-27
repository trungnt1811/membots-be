package fair_launch

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
)

const keyPrefixFairLauch = "fair_lauch_"
const cacheTimeFairLauch = 3 * time.Second

type fairLauchCache struct {
	FairLauchRepository		interfaces.FairLauchRepository
	Cache                   caching.Repository
}

func NewFairLauchCacheRepository(repo interfaces.FairLauchRepository,
	cache caching.Repository,
) interfaces.FairLauchRepository {
	return &fairLauchCache{
		FairLauchRepository: repo,
		Cache:                     cache,
	}
}

func (c fairLauchCache) GetMeme20MetaByTicker(ctx context.Context, ticker string) (model.Meme20Meta, error) {
	key := &caching.Keyer{Raw: keyPrefixFairLauch + fmt.Sprint("GetMeme20MetaByTicker_", ticker)}
	var meme20Meta model.Meme20Meta
	err := c.Cache.RetrieveItem(key, &meme20Meta)
	if err != nil {
		// cache miss
		meme20Meta, err = c.GetMeme20MetaByTicker(ctx, ticker)
		if err != nil {
			return meme20Meta, err
		}
		if err = c.Cache.SaveItem(key, meme20Meta, cacheTimeFairLauch); err != nil {
			return meme20Meta, err
		}
	}
	return meme20Meta, nil
}

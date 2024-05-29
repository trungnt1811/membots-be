package memeception

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
)

const (
	keyPrefixMemeception = "memeception_"
	cacheTimeMemeception = 10 * time.Second
)

type memeceptionCache struct {
	MemeceptionRepository interfaces.MemeceptionRepository
	Cache                 caching.Repository
}

func NewMemeceptionCacheRepository(repo interfaces.MemeceptionRepository,
	cache caching.Repository,
) interfaces.MemeceptionRepository {
	return &memeceptionCache{
		MemeceptionRepository: repo,
		Cache:                 cache,
	}
}

func (c memeceptionCache) GetMemeceptionBySymbol(ctx context.Context, symbol string) (model.Meme, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint("GetMemeceptionBySymbol_", symbol)}
	var memeMeta model.Meme
	err := c.Cache.RetrieveItem(key, &memeMeta)
	if err != nil {
		// cache miss
		memeMeta, err = c.GetMemeceptionBySymbol(ctx, symbol)
		if err != nil {
			return memeMeta, err
		}
		if err = c.Cache.SaveItem(key, memeMeta, cacheTimeMemeception); err != nil {
			return memeMeta, err
		}
	}
	return memeMeta, nil
}

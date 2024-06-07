package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

const (
	keyPrefixStats = "stats_"
	cacheTimeStats = 5 * time.Second
)

type statsCache struct {
	StatsUCase interfaces.StatsUCase
	Cache      caching.Repository
}

func NewStatsCacheUCase(repo interfaces.StatsUCase,
	cache caching.Repository,
) interfaces.StatsUCase {
	return &statsCache{
		StatsUCase: repo,
		Cache:      cache,
	}
}

func (s *statsCache) GetStatsByMemeAddress(ctx context.Context, address string) (interface{}, error) {
	key := &caching.Keyer{Raw: keyPrefixStats + fmt.Sprint("GetStatsByMemeAddress_", address)}
	var swaps interface{}
	err := s.Cache.RetrieveItem(key, &swaps)
	if err != nil {
		// cache miss
		swaps, err = s.StatsUCase.GetStatsByMemeAddress(ctx, address)
		if err != nil {
			return swaps, err
		}
		if err = s.Cache.SaveItem(key, swaps, cacheTimeStats); err != nil {
			return swaps, err
		}
	}
	return swaps, nil
}

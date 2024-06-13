package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
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

func (s *statsCache) GetStatsByMemeAddress(ctx context.Context, address string) (dto.TokenStatsResp, error) {
	key := &caching.Keyer{Raw: keyPrefixStats + fmt.Sprint("GetStatsByMemeAddress_", address)}
	var stats dto.TokenStatsResp
	err := s.Cache.RetrieveItem(key, &stats)
	if err != nil {
		// cache miss
		stats, err = s.StatsUCase.GetStatsByMemeAddress(ctx, address)
		if err != nil {
			return stats, err
		}
		if err = s.Cache.SaveItem(key, stats, cacheTimeStats); err != nil {
			return stats, err
		}
	}
	return stats, nil
}

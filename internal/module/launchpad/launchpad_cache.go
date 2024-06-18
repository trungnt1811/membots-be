package launchpad

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

const (
	keyPrefixLaunchpad = "launchpad_"
	cacheTimeLaunchpad = 5 * time.Second
)

type launchpadCache struct {
	LaunchpadUCase interfaces.LaunchpadUCase
	Cache          caching.Repository
}

func NewLaunchpadCacheUCase(repo interfaces.LaunchpadUCase,
	cache caching.Repository,
) interfaces.LaunchpadUCase {
	return &launchpadCache{
		LaunchpadUCase: repo,
		Cache:          cache,
	}
}

func (l *launchpadCache) GetHistory(ctx context.Context, address string) (dto.LaunchpadInfoResp, error) {
	key := &caching.Keyer{Raw: keyPrefixLaunchpad + fmt.Sprint("GetHistory_", address)}
	var launchpadInfoRsp dto.LaunchpadInfoResp
	err := l.Cache.RetrieveItem(key, &launchpadInfoRsp)
	if err != nil {
		// cache miss
		launchpadInfoRsp, err = l.LaunchpadUCase.GetHistory(ctx, address)
		if err != nil {
			return launchpadInfoRsp, err
		}
		if err = l.Cache.SaveItem(key, launchpadInfoRsp, cacheTimeLaunchpad); err != nil {
			return launchpadInfoRsp, err
		}
	}
	return launchpadInfoRsp, nil
}

package swap

import (
	"context"
	"fmt"
	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"time"
)

const (
	keyPrefixSwap = "swap_"
	cacheTimeSwap = 5 * time.Second
)

type swapCache struct {
	SwapUCase interfaces.SwapUCase
	Cache     caching.Repository
}

func NewSwapCacheRepository(repo interfaces.SwapUCase,
	cache caching.Repository,
) interfaces.SwapUCase {
	return &swapCache{
		SwapUCase: repo,
		Cache:     cache,
	}
}

func (s *swapCache) GetSwaps(ctx context.Context, address string) (dto.SwapHistoryByAddressRsp, error) {
	key := &caching.Keyer{Raw: keyPrefixSwap + fmt.Sprint("GetSwaps_", address)}
	var swaps dto.SwapHistoryByAddressRsp
	err := s.Cache.RetrieveItem(key, &swaps)
	if err != nil {
		// cache miss
		swaps, err = s.GetSwaps(ctx, address)
		if err != nil {
			return swaps, err
		}
		if err = s.Cache.SaveItem(key, swaps, cacheTimeSwap); err != nil {
			return swaps, err
		}
	}
	return swaps, nil
}

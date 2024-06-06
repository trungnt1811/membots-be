package swap

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

const (
	keyPrefixSwap = "swap_"
	cacheTimeSwap = 5 * time.Second
)

type swapCache struct {
	SwapUCase interfaces.SwapUCase
	Cache     caching.Repository
}

func (s *swapCache) GetQuote(ctx context.Context, url string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func NewSwapCacheUCase(repo interfaces.SwapUCase,
	cache caching.Repository,
) interfaces.SwapUCase {
	return &swapCache{
		SwapUCase: repo,
		Cache:     cache,
	}
}

func (s *swapCache) GetSwaps(ctx context.Context, address string) (dto.SwapHistoryByAddressResp, error) {
	key := &caching.Keyer{Raw: keyPrefixSwap + fmt.Sprint("GetSwaps_", address)}
	var swaps dto.SwapHistoryByAddressResp
	err := s.Cache.RetrieveItem(key, &swaps)
	if err != nil {
		// cache miss
		swaps, err = s.SwapUCase.GetSwaps(ctx, address)
		if err != nil {
			return swaps, err
		}
		if err = s.Cache.SaveItem(key, swaps, cacheTimeSwap); err != nil {
			return swaps, err
		}
	}
	return swaps, nil
}

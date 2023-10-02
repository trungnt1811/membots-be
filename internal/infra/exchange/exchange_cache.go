package exchange

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

const keyPrefix = "tiki_client_"

type tikiClientCache struct {
	TikiClient interfaces.TokenPriceRepo
	Client     caching.RedisClient
}

func NewTikiClientCache(tikiClient interfaces.TokenPriceRepo,
	client caching.RedisClient) interfaces.TokenPriceRepo {
	return tikiClientCache{
		TikiClient: tikiClient,
		Client:     client,
	}
}

func (c tikiClientCache) GetAstraPrice(ctx context.Context) (int64, error) {
	key := keyPrefix + "GetAstraPrice"
	var price int64
	err := c.Client.Get(ctx, key, &price)
	if err != nil {
		// cache miss
		price, err = c.TikiClient.GetAstraPrice(ctx)
		if err != nil {
			return price, err
		}
		if err = c.Client.Set(ctx, key, price, 2*time.Minute); err != nil {
			return price, err
		}
	}
	return price, nil
}

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
	cacheTimeMemeception = 5 * time.Second
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

func (c memeceptionCache) CreateMeme(ctx context.Context, model model.Meme) error {
	return c.MemeceptionRepository.CreateMeme(ctx, model)
}

func (c memeceptionCache) UpdateMeme(ctx context.Context, model model.Meme) error {
	return c.MemeceptionRepository.UpdateMeme(ctx, model)
}

func (c memeceptionCache) UpdateMemeception(ctx context.Context, model model.Memeception) error {
	return c.MemeceptionRepository.UpdateMemeception(ctx, model)
}

func (c memeceptionCache) GetListMemeProcessing(ctx context.Context) ([]model.MemeOnchainInfo, error) {
	// TODO: implement later
	return nil, nil
}

func (c memeceptionCache) GetListMemeLive(ctx context.Context) ([]model.MemeOnchainInfo, error) {
	// TODO: implement later
	return nil, nil
}

func (c memeceptionCache) GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint("GetMemeceptionByContractAddress_", contractAddress)}
	var memeMeta model.Meme
	err := c.Cache.RetrieveItem(key, &memeMeta)
	if err != nil {
		// cache miss
		memeMeta, err = c.MemeceptionRepository.GetMemeceptionByContractAddress(ctx, contractAddress)
		if err != nil {
			return memeMeta, err
		}
		if err = c.Cache.SaveItem(key, memeMeta, cacheTimeMemeception); err != nil {
			return memeMeta, err
		}
	}
	return memeMeta, nil
}

func (c memeceptionCache) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	// TODO: implement later
	return nil, nil
}

func (c memeceptionCache) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	// TODO: implement later
	return nil, nil
}

func (c memeceptionCache) GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error) {
	// TODO: implement later
	return nil, nil
}

func (c memeceptionCache) GetMapMemeSymbolAndLogoURL(ctx context.Context, contractAddresses []string) (
	map[string]model.MemeSymbolAndLogoURL,
	error,
) {
	// TODO: implement later
	return nil, nil
}

func (c memeceptionCache) GetMemeceptionBySymbol(ctx context.Context, symbol string) (model.Meme, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint("GetMemeceptionBySymbol_", symbol)}
	var meme model.Meme
	err := c.Cache.RetrieveItem(key, &meme)
	if err != nil {
		// cache miss
		meme, err = c.MemeceptionRepository.GetMemeceptionBySymbol(ctx, symbol)
		if err != nil {
			return meme, err
		}
		if err = c.Cache.SaveItem(key, meme, cacheTimeMemeception); err != nil {
			return meme, err
		}
	}
	return meme, nil
}

func (c memeceptionCache) GetMemeIDAndStartAtByContractAddress(ctx context.Context, contractAddress string) (
	model.MemeceptionMemeIDAndStartAt,
	error,
) {
	// TODO: implement later
	return model.MemeceptionMemeIDAndStartAt{}, nil
}

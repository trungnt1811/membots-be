package memeception

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flexstack.ai/membots-be/internal/constant"
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
	key := &caching.Keyer{Raw: keyPrefixMemeception + "GetListMemeProcessing"}
	var memes []model.MemeOnchainInfo
	err := c.Cache.RetrieveItem(key, &memes)
	if err != nil {
		// cache miss
		memes, err = c.MemeceptionRepository.GetListMemeProcessing(ctx)
		if err != nil {
			return memes, err
		}
		if err = c.Cache.SaveItem(key, memes, cacheTimeMemeception); err != nil {
			return memes, err
		}
	}
	return memes, nil
}

func (c memeceptionCache) GetListMemeLive(ctx context.Context) ([]model.MemeOnchainInfo, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + "GetListMemeLive"}
	var memes []model.MemeOnchainInfo
	err := c.Cache.RetrieveItem(key, &memes)
	if err != nil {
		// cache miss
		memes, err = c.MemeceptionRepository.GetListMemeLive(ctx)
		if err != nil {
			return memes, err
		}
		if err = c.Cache.SaveItem(key, memes, cacheTimeMemeception); err != nil {
			return memes, err
		}
	}
	return memes, nil
}

func (c memeceptionCache) GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint("GetMemeceptionByContractAddress_", contractAddress)}
	var meme model.Meme
	err := c.Cache.RetrieveItem(key, &meme)
	if err != nil {
		// cache miss
		meme, err = c.MemeceptionRepository.GetMemeceptionByContractAddress(ctx, contractAddress)
		if err != nil {
			return meme, err
		}
		if err = c.Cache.SaveItem(key, meme, cacheTimeMemeception); err != nil {
			return meme, err
		}
	}
	return meme, nil
}

func (c memeceptionCache) GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + "GetMemeceptionsPast"}
	var memeceptions []model.Memeception
	err := c.Cache.RetrieveItem(key, &memeceptions)
	if err != nil {
		// cache miss
		memeceptions, err = c.MemeceptionRepository.GetMemeceptionsPast(ctx)
		if err != nil {
			return memeceptions, err
		}
		if err = c.Cache.SaveItem(key, memeceptions, cacheTimeMemeception); err != nil {
			return memeceptions, err
		}
	}
	return memeceptions, nil
}

func (c memeceptionCache) GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + "GetMemeceptionsLive"}
	var memeceptions []model.Memeception
	err := c.Cache.RetrieveItem(key, &memeceptions)
	if err != nil {
		// cache miss
		memeceptions, err = c.MemeceptionRepository.GetMemeceptionsLive(ctx)
		if err != nil {
			return memeceptions, err
		}
		if err = c.Cache.SaveItem(key, memeceptions, cacheTimeMemeception); err != nil {
			return memeceptions, err
		}
	}
	return memeceptions, nil
}

func (c memeceptionCache) GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + "GetMemeceptionsLatest"}
	var memeceptions []model.Memeception
	err := c.Cache.RetrieveItem(key, &memeceptions)
	if err != nil {
		// cache miss
		memeceptions, err = c.MemeceptionRepository.GetMemeceptionsLatest(ctx)
		if err != nil {
			return memeceptions, err
		}
		if err = c.Cache.SaveItem(key, memeceptions, cacheTimeMemeception); err != nil {
			return memeceptions, err
		}
	}
	return memeceptions, nil
}

func (c memeceptionCache) GetMapMemeSymbolAndLogoURL(ctx context.Context, contractAddresses []string) (
	map[string]model.MemeSymbolAndLogoURL,
	error,
) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint(
		"GetMapMemeSymbolAndLogoURL_", strings.Join(contractAddresses, "_"),
	)}
	var mapMeme map[string]model.MemeSymbolAndLogoURL
	err := c.Cache.RetrieveItem(key, &mapMeme)
	if err != nil {
		// cache miss
		mapMeme, err = c.MemeceptionRepository.GetMapMemeSymbolAndLogoURL(ctx, contractAddresses)
		if err != nil {
			return mapMeme, err
		}
		if err = c.Cache.SaveItem(key, mapMeme, cacheTimeMemeception); err != nil {
			return mapMeme, err
		}
	}
	return mapMeme, nil
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
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint(
		"GetMemeIDAndStartAtByContractAddress_", contractAddress,
	)}
	var meme model.MemeceptionMemeIDAndStartAt
	err := c.Cache.RetrieveItem(key, &meme)
	if err != nil {
		// cache miss
		meme, err = c.MemeceptionRepository.GetMemeIDAndStartAtByContractAddress(ctx, contractAddress)
		if err != nil {
			return meme, err
		}
		if err = c.Cache.SaveItem(key, meme, cacheTimeMemeception); err != nil {
			return meme, err
		}
	}
	return meme, nil
}

func (c memeceptionCache) MemeceptionExists(ctx context.Context, symbol string) (bool, error) {
	key := &caching.Keyer{Raw: keyPrefixMemeception + fmt.Sprint("MemeceptionExists_", symbol)}
	var isExist bool
	err := c.Cache.RetrieveItem(key, &isExist)
	if err != nil {
		// cache miss
		isExist, err = c.MemeceptionRepository.MemeceptionExists(ctx, symbol)
		if err != nil {
			return isExist, err
		}
		if err = c.Cache.SaveItem(key, isExist, cacheTimeMemeception); err != nil {
			return isExist, err
		}
	}
	return isExist, nil
}

func (c memeceptionCache) GetETHPrice() (uint64, error) {
	key := &caching.Keyer{Raw: "GetETHPrice"}
	var price uint64
	err := c.Cache.RetrieveItem(key, &price)
	if err != nil {
		// cache miss
		price, err = c.MemeceptionRepository.GetETHPrice()
		if err != nil {
			return price, err
		}
		if err = c.Cache.SaveItem(key, price, constant.CacheTimeETHPrice); err != nil {
			return price, err
		}
	}
	return price, nil
}

package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/vmihailenco/msgpack/v5"
)

type Repository interface {
	SaveItem(key fmt.Stringer, val interface{}, expire time.Duration) error
	RetrieveItem(key fmt.Stringer, val interface{}) error
	ExpireItem(key fmt.Stringer) error
	RemoveItem(key fmt.Stringer) error
}

type CachingRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewCachingRepository(ctx context.Context, client *redis.Client) *CachingRepository {
	return &CachingRepository{
		ctx:    ctx,
		client: client,
	}
}

func (repo *CachingRepository) SaveItem(key fmt.Stringer, val interface{}, expire time.Duration) error {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}
	status := repo.client.Set(repo.ctx, key.String(), b, expire)
	err = status.Err()
	if err != nil {
		return err
	}
	return nil
}

func (repo *CachingRepository) RetrieveItem(key fmt.Stringer, val interface{}) error {
	status := repo.client.Get(repo.ctx, key.String())
	res, err := status.Bytes()
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal(res, val)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CachingRepository) ExpireItem(key fmt.Stringer) error {
	status := repo.client.Expire(repo.ctx, key.String(), time.Millisecond*100)
	return status.Err()
}

func (repo *CachingRepository) RemoveItem(key fmt.Stringer) error {
	status := repo.client.Del(repo.ctx, key.String())
	return status.Err()
}

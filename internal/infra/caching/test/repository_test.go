package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/astraprotocol/membots-be/internal/infra/caching"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type Keyer struct {
	Raw string
}

type Value struct {
	Raw string
}

func (k *Keyer) String() string {
	return k.Raw
}

func getRdb() *redis.Client {
	addr := "localhost:6379"
	if os.Getenv("ENV_FILE") == "test.env" {
		addr = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	return rdb
}

func TestSaveRetrieveRemoveItem(t *testing.T) {
	asserts := assert.New(t)
	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := &Keyer{Raw: "test_key"}
	val := &Value{Raw: "hello world"}

	rdb := getRdb()

	repo := caching.NewCachingRepository(testCtx, rdb)

	err := repo.SaveItem(key, val, time.Duration(0))

	asserts.Nil(err, "Save item error!")

	var newVal Value
	err = repo.RetrieveItem(key, &newVal)

	asserts.Nil(err, "Retrieve item error!")
	// asserts.NotNil(newVal, "Retrieve item must not be nil!")

	asserts.Equal(val.Raw, newVal.Raw, "Retrieve item not equal")

	err = repo.RemoveItem(key)

	asserts.Nil(err, "Remove item error")
}

func TestSaveItemWithExpire(t *testing.T) {
	asserts := assert.New(t)
	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := &Keyer{Raw: "test_key"}
	val := &Value{Raw: "hello world"}

	rdb := getRdb()
	repo := caching.NewCachingRepository(testCtx, rdb)

	duration := time.Millisecond * 500
	// caching for a duration
	err := repo.SaveItem(key, val, duration)
	asserts.Nil(err, "Save item error!")

	// Expect retrieve item before expired
	var newVal Value
	err = repo.RetrieveItem(key, &newVal)
	asserts.Nil(err, "Retrieve item error!")
	// asserts.NotNil(newVal, "Retrieve item must not be nil!")
	asserts.Equal(val.Raw, newVal.Raw, "Retrieve item not equal")

	// Expect retrieve item failed after expired
	time.Sleep(duration)
	var emptyVal Value
	err = repo.RetrieveItem(key, &emptyVal)
	asserts.NotNil(err, "Retrieve item error")
	asserts.Equal(emptyVal.Raw, "", "Retrieve item must not available")
}

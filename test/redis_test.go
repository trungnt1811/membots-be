package test

import (
	"context"
	"testing"
	"time"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/stretchr/testify/assert"
)

func TestRedisConn(t *testing.T) {
	asserts := assert.New(t)
	rdc := conf.RedisConn()
	deadline, ok := t.Deadline()
	var testCtx context.Context
	if !ok {
		testCtx = context.Background()
	} else {
		testCtx, _ = context.WithDeadline(context.Background(), deadline)
	}

	cmd := rdc.Set(testCtx, "test", "test", time.Duration(time.Second*30))
	err := cmd.Err()
	asserts.NoError(err)

	strCmd := rdc.Get(testCtx, "test")
	asserts.NoError(strCmd.Err())
	asserts.Equal("test", strCmd.Val())
}

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/flexstack.ai/membots-be/conf"
)

func TestRedisConn(t *testing.T) {
	asserts := assert.New(t)
	rdc := conf.RedisConn()
	deadline, ok := t.Deadline()
	var testCtx context.Context
	var cancelFunc context.CancelFunc
	if !ok {
		testCtx = context.Background()
	} else {
		testCtx, cancelFunc = context.WithDeadline(context.Background(), deadline)
		fmt.Printf("cancelFunc: %v\n", cancelFunc)
	}

	cmd := rdc.Set(testCtx, "test", "test", time.Second*30)
	err := cmd.Err()
	asserts.NoError(err)

	strCmd := rdc.Get(testCtx, "test")
	asserts.NoError(strCmd.Err())
	asserts.Equal("test", strCmd.Val())
}

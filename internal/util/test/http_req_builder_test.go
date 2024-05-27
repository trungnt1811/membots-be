package test

import (
	"testing"

	"github.com/flexstack.ai/membots-be/internal/util"
	"github.com/stretchr/testify/assert"
)

func Test_ReqBuilder(t *testing.T) {
	asserts := assert.New(t)

	// Start an echo http server at 8080 and run test
	resp, err := util.NewHttpRequestBuilder().SetBody(map[string]any{
		"content": "Test HTTP API",
	}).Build().Post("http://localhost:8080")

	asserts.NoError(err)
	asserts.True(resp.IsSuccess())
}

package test

import (
	"fmt"
	"testing"

	"github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/stretchr/testify/assert"
)

func Test_PackQueryParamsToNoSearchUrl(t *testing.T) {
	asserts := assert.New(t)

	url := "https://astranaut.io"
	params := map[string]string{
		"utm_source": "stella",
	}

	newUrl := util.PackQueryParamsToUrl(url, params)

	asserts.Equal("https://astranaut.io?utm_source=stella", newUrl)
}

func Test_PackQueryParamsToFullUrl(t *testing.T) {
	asserts := assert.New(t)

	url1 := "https://astranaut.io?q=1"
	params := map[string]string{
		"utm_source": "stella",
	}

	newUrl := util.PackQueryParamsToUrl(url1, params)

	asserts.Equal("https://astranaut.io?q=1&utm_source=stella", newUrl)

	url2 := "https://astranaut.io"
	asserts.Equal("https://astranaut.io?utm_source=stella", util.PackQueryParamsToUrl(url2, params))
}

func Test_ParseAndStringifyUTMContent(t *testing.T) {
	asserts := assert.New(t)

	userId := uint(1)
	trackedId := uint64(100)
	utmStr := fmt.Sprintf("%d-%d", userId, trackedId)

	asserts.Equal(utmStr, util.StringifyUTMContent(userId, trackedId))
	parsedUserId, parsedTrackedId := util.ParseUTMContent(utmStr)
	asserts.Equal(userId, parsedUserId)
	asserts.Equal(trackedId, parsedTrackedId)
}

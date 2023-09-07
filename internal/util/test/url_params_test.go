package test

import (
	"github.com/astraprotocol/affiliate-system/internal/util"
	"testing"

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

	url := "https://astranaut.io?q=1"
	params := map[string]string{
		"utm_source": "stella",
	}

	newUrl := util.PackQueryParamsToUrl(url, params)

	asserts.Equal("https://astranaut.io?q=1&utm_source=stella", newUrl)
}

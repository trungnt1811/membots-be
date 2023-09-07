package util

import (
	"fmt"
	"net/url"
	"strings"
)

func PackQueryParamsToUrl(s string, params map[string]string) string {
	obj, err := url.Parse(s)
	if err != nil {
		return s
	}

	// Pack new params
	elems := []string{}
	if obj.RawQuery != "" {
		elems = append(elems, obj.RawQuery)
	}
	for key, val := range params {
		elems = append(elems, fmt.Sprintf("%s=%s", key, url.QueryEscape(val)))
	}
	obj.RawQuery = strings.Join(elems, "&")

	return obj.String()
}

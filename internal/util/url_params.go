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

	q := []string{}
	for key, val := range params {
		q = append(q, fmt.Sprint(key, "=", val))
	}

	if obj.RawQuery == "" {
		return fmt.Sprint(s, "?", strings.Join(q, "&"))
	}

	return fmt.Sprint(s, "&", strings.Join(q, "&"))
}

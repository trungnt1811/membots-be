package util

import (
	b64 "encoding/base64"
	"fmt"
	"net/url"
	"strconv"
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

// The function `ParseUTMContent` takes a base64 encoded string as input, decodes it, splits it into
// parts, and returns the user ID and tracked ID as unsigned integers.
func ParseUTMContent(utmContent string) (uint, uint64) {
	decoded, err := b64.URLEncoding.DecodeString(utmContent)
	if err != nil {
		// Cannot decode utm, return 0
		return 0, 0
	}
	parts := strings.Split(string(decoded), "-")
	if len(parts) == 0 {
		return 0, 0
	}

	userId, _ := strconv.ParseUint(parts[0], 10, 64)
	if len(parts) == 1 {
		return uint(userId), 0
	}
	trackedId, _ := strconv.ParseUint(parts[1], 10, 64)

	return uint(userId), trackedId
}

// The function takes a user ID and a tracked ID, combines them into a string, encodes the string using
// base64 encoding, and returns the encoded string.
func StringifyUTMContent(userId uint, trackedId uint64) string {
	s := fmt.Sprintf("%d-%d", userId, trackedId)
	encoded := b64.URLEncoding.EncodeToString([]byte(s))

	return encoded
}

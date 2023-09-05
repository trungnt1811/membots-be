package test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade/types"
	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
)

func TestParseOrderQuery(t *testing.T) {
	asserts := assert.New(t)
	sinceStr := "2023-08-22T00:00:00Z"
	untilStr := "2023-08-21T00:00:00Z"

	since, _ := time.Parse(time.RFC3339, sinceStr)
	until, _ := time.Parse(time.RFC3339, untilStr)

	empty := types.ATOrderQuery{
		Since: since,
		Until: until,
	}

	v, err := query.Values(&empty)
	asserts.NoError(err)
	expected := fmt.Sprintf("since=%s&until=%s", url.QueryEscape(sinceStr), url.QueryEscape(untilStr))
	asserts.Equal(expected, v.Encode())
}

func TestCustomTime(t *testing.T) {
	asserts := assert.New(t)
	dStr := "2022-08-08T22:54:03"
	bytes := []byte(fmt.Sprintf("{\"date\": \"%s\"}", dStr))
	var res struct {
		Date types.CustomTime `json:"date"`
	}
	err := json.Unmarshal(bytes, &res)
	asserts.NoError(err)

	bytes, err = json.Marshal(&res)
	asserts.NoError(err)
	asserts.Contains(string(bytes), dStr)
}

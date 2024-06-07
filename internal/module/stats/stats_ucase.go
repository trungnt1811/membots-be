package stats

import (
	"context"
	"encoding/json"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/util"
)

type statsUCase struct{}

func NewStatsUcase() interfaces.StatsUCase {
	return &statsUCase{}
}

func (uc *statsUCase) GetStatsByMemeAddress(ctx context.Context, url string) (interface{}, error) {
	resp, err := util.NewHttpRequestBuilder().Build().Get(url)
	if err != nil {
		return nil, err
	}
	var result interface{} // This can be any type that you expect the JSON to conform to
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err // Handle JSON parsing error
	}

	return result, nil
}

package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/util"
)

type statsUCase struct {
	Client                *subgraphclient.Client
	MemeceptionRepository interfaces.MemeceptionRepository
}

func NewStatsUcase(client *subgraphclient.Client, memeceptionRepository interfaces.MemeceptionRepository) interfaces.StatsUCase {
	return &statsUCase{Client: client, MemeceptionRepository: memeceptionRepository}
}

func (uc *statsUCase) GetStatsByMemeAddress(ctx context.Context, address string) (dto.TokenStatsResp, error) {
	meme, err := uc.MemeceptionRepository.GetMemeIDAndStartAtByContractAddress(ctx, address)
	if err != nil {
		return dto.TokenStatsResp{}, err
	}

	requestOpts := &subgraphclient.RequestOptions{
		IncludeFields: []string{
			"id",
			"txCount",
			"volumeUSD",
			"totalValueLockedUSD",
			"tokenDayData.date",
			"tokenDayData.priceUSD",
			"tokenDayData.close",
			"tokenDayData.volumeUSD",
			"tokenDayData.totalValueLockedUSD",
		},
	}

	// TODO(perf): query by both token address and date
	response, err := uc.Client.GetTokenById(ctx, address, requestOpts)
	if err != nil {
		return dto.TokenStatsResp{}, err
	}

	// Get today data
	now := time.Now()
	midnight := now.Truncate(24 * time.Hour)
	tokenDayData := subgraphclient.TokenDayDataET{}
	for _, tdd := range response.Token.TokenDayData {
		if tdd.Date == uint64(midnight.Unix()) {
			tokenDayData.VolumeUSD = tdd.VolumeUSD
			tokenDayData.Close = tdd.Close
			tokenDayData.TotalValueLockedUSD = tdd.TotalValueLockedUSD
			break
		}
	}

	resp := dto.TokenStatsResp{
		Stats: dto.Stats{
			MemeID:              fmt.Sprint(meme.MemeID),
			StartDate:           time.Unix(int64(meme.StartAt), 0).UTC().Format(time.RFC3339Nano),
			VolumeUSD:           response.Token.VolumeUSD,
			SwapCount:           util.StringToUint64(response.Token.TxCount),
			CloseUSD:            tokenDayData.Close,
			TotalValueLockedUSD: response.Token.TotalValueLockedUSD,
			HoldersCount:        0, // TODO: implement get token holders
		},
		DailyDiffPct: dto.DailyDiffPct{
			VolumeUSD:           util.StringToUint64(tokenDayData.VolumeUSD),
			SwapCount:           0, // TODO: implement get count total tx by day
			CloseUSD:            util.StringToUint64(tokenDayData.Close),
			TotalValueLockedUSD: util.StringToUint64(tokenDayData.TotalValueLockedUSD),
			HoldersCount:        0, // TODO: implement get token holders by day
		},
	}
	return resp, nil
}

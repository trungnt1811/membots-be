package swap

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/flexstack.ai/membots-be/internal/util"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type swapUCase struct {
	Client *subgraphclient.Client
}

func NewSwapUcase(client *subgraphclient.Client) interfaces.SwapUCase {
	return &swapUCase{Client: client}
}

func (uc *swapUCase) GetSwaps(ctx context.Context, address string) (dto.SwapHistoryByAddressResp, error) {
	requestOpts := &subgraphclient.RequestOptions{
		IncludeFields: []string{
			"id",
			"timestamp",
			"amount0",
			"amount1",
			"amountUSD",
			"sqrtPriceX96",
			"sender",
			"recipient",
			"transaction.id",
		},
	}

	response, err := uc.Client.GetSwapHistoryByPoolId(ctx, address, requestOpts)
	if err != nil {
		return dto.SwapHistoryByAddressResp{}, err
	}
	// convert response to dto.SwapHistoryByAddressRsp
	var swaps []dto.Swap
	for _, swap := range response.Swaps {
		timestamp, err := strconv.ParseUint(swap.Timestamp, 10, 64)
		if err != nil {
			timestamp = 0
		}
		convertedSwap := dto.Swap{
			Buy:           swap.Amount0 > "0",
			Amount0:       swap.Amount0,
			Amount1:       swap.Amount1,
			AmountUSD:     swap.AmountUSD,
			PriceETH:      swap.SqrtPriceX96,
			PriceUSD:      "",
			WalletAddress: swap.Sender,
			TxHash:        swap.Transaction.ID,
			SwapAt:        timestamp,
			Token1IsMeme:  true,
		}
		swaps = append(swaps, convertedSwap)
	}
	swapHistoryByAddressRsp := dto.SwapHistoryByAddressResp{
		Swaps: swaps,
	}
	return swapHistoryByAddressRsp, nil
}

func (uc *swapUCase) GetQuote(ctx context.Context, url string) (interface{}, error) {
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

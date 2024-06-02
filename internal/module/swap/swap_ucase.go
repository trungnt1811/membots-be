package swap

import (
	"context"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"strconv"

	unigraphclient "github.com/emersonmacro/go-uniswap-subgraph-client"
	"github.com/flexstack.ai/membots-be/internal/dto"
)

type swapUCase struct {
	Client *unigraphclient.Client
}

func NewSwapUcase(client *unigraphclient.Client) interfaces.SwapUCase {
	return &swapUCase{Client: client}
}

func (uc *swapUCase) GetSwaps(ctx context.Context, address string) (dto.SwapHistoryByAddressRsp, error) {
	requestOpts := &unigraphclient.RequestOptions{
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
		return dto.SwapHistoryByAddressRsp{}, err
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
	swapHistoryByAddressRsp := dto.SwapHistoryByAddressRsp{
		Swaps: swaps,
	}
	return swapHistoryByAddressRsp, nil
}

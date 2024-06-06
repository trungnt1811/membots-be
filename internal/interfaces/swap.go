package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
)

type SwapUCase interface {
	GetSwaps(ctx context.Context, address string) (dto.SwapHistoryByAddressResp, error)
	GetQuote(ctx context.Context, url string) (interface{}, error)
}

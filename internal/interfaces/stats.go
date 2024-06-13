package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
)

type StatsUCase interface {
	GetStatsByMemeAddress(ctx context.Context, address string) (dto.TokenStatsResp, error)
}

package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
)

type LaunchpadUCase interface {
	GetHistory(ctx context.Context, address string) (dto.LaunchpadInfoResp, error)
}

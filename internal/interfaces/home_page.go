package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

type HomePageUCase interface {
	GetHomePage(ctx context.Context, userId uint64) (dto.HomePageDto, error)
}

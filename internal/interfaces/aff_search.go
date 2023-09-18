package interfaces

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffSearchUCase interface {
	Search(ctx context.Context, q string, page, size int) (dto.AffSearchResponseDto, error)
	SearchWithStatus(ctx context.Context, q string, stellaStatus []string, page, size int) (dto.AffSearchResponseDto, error)
}

type AffSearchRepository interface {
	Search(ctx context.Context, q string, page, size int) (model.AffSearch, error)
	SearchWithStatus(ctx context.Context, q string, stellaStatus []string, page, size int) (model.AffSearch, error)
}

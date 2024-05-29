package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type MemeceptionRepository interface {
	GetMemeceptionBySymbol(ctx context.Context, symbol string) (model.Meme, error)
}

type MemeceptionUCase interface {
	GetMemeceptionBySymbol(ctx context.Context, symbol string) (dto.MemeceptionResp, error)
}

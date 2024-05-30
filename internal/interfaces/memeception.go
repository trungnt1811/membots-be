package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type MemeceptionRepository interface {
	GetMemeceptionBySymbol(ctx context.Context, symbol string) (model.Meme, error)
	GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error)
	GetMemeceptionsUpcoming(ctx context.Context) ([]model.Memeception, error)
	GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error)
}

type MemeceptionUCase interface {
	GetMemeceptionBySymbol(ctx context.Context, symbol string) (dto.MemeceptionDetailResp, error)
	GetMemeceptions(ctx context.Context) (dto.MemeceptionResp, error)
}

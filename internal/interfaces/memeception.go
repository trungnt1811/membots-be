package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type MemeceptionRepository interface {
	CreateMeme(ctx context.Context, model model.Meme) error
	GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error)
	GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error)
	GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error)
	GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error)
}

type MemeceptionUCase interface {
	CreateMeme(ctx context.Context, payload dto.CreateMemePayload) error
	GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (dto.MemeceptionDetailResp, error)
	GetMemeceptions(ctx context.Context) (dto.MemeceptionsResp, error)
}

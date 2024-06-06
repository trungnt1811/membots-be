package interfaces

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type MemeceptionRepository interface {
	CreateMeme(ctx context.Context, model model.Meme) error
	UpdateMeme(ctx context.Context, model model.Meme) error
	UpdateMemeception(ctx context.Context, model model.Memeception) error
	GetListMemeProcessing(ctx context.Context) ([]model.MemeOnchainInfo, error)
	GetListMemeLive(ctx context.Context) ([]model.MemeOnchainInfo, error)
	GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (model.Meme, error)
	GetMemeceptionsPast(ctx context.Context) ([]model.Memeception, error)
	GetMemeceptionsLive(ctx context.Context) ([]model.Memeception, error)
	GetMemeceptionsLatest(ctx context.Context) ([]model.Memeception, error)
	GetMapMemeSymbolAndLogoURL(ctx context.Context, contractAddresses []string) (map[string]model.MemeSymbolAndLogoURL, error)
}

type MemeceptionUCase interface {
	CreateMeme(ctx context.Context, payload dto.CreateMemePayload) error
	GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (dto.MemeceptionDetailResp, error)
	GetMemeceptions(ctx context.Context) (dto.MemeceptionsResp, error)
}

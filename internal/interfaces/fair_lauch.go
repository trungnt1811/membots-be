package interfaces

import (
	"context"

	"github.com/astraprotocol/membots-be/internal/dto"
	"github.com/astraprotocol/membots-be/internal/model"
)

type FairLauchRepository interface {
	GetMeme20MetaByTicker(ctx context.Context, ticker string) (model.Meme20Meta, error)
}

type FairLauchUCase interface {
	GetMeme20MetaByTicker(ctx context.Context, ticker string) (dto.Meme20MetaDto, error)
}

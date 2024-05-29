package memeception

import (
	"context"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type memeceptionUCase struct {
	MemeceptionRepository interfaces.MemeceptionRepository
}

func NewMemeceptionUCase(memeceptionRepository interfaces.MemeceptionRepository) interfaces.MemeceptionUCase {
	return &memeceptionUCase{
		MemeceptionRepository: memeceptionRepository,
	}
}

func (u *memeceptionUCase) GetMemeceptionBySymbol(ctx context.Context, symbol string) (dto.MemeceptionResp, error) {
	memeMeta, err := u.MemeceptionRepository.GetMemeceptionBySymbol(ctx, symbol)
	if err != nil {
		return dto.MemeceptionResp{}, err
	}
	// TODO: get price from RPC
	memeceptionResp := dto.MemeceptionResp{
		Meme:  memeMeta.ToDto(),
		Price: 0,
	}
	return memeceptionResp, nil
}

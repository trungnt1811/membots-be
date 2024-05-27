package fair_launch

import (
	"context"

	"github.com/astraprotocol/membots-be/internal/dto"
	"github.com/astraprotocol/membots-be/internal/interfaces"
)

type fairLauchUCase struct {
	FairLauchRepository interfaces.FairLauchRepository
}

func NewFairLauchUCaseUCase(fairLauchRepository interfaces.FairLauchRepository) interfaces.FairLauchUCase {
	return &fairLauchUCase{
		FairLauchRepository: fairLauchRepository,
	}
}

func (u *fairLauchUCase) GetMeme20MetaByTicker(ctx context.Context, ticker string) (dto.Meme20MetaDto, error) {
	meme20Meta, err := u.FairLauchRepository.GetMeme20MetaByTicker(ctx, ticker)
	if err != nil {
		return dto.Meme20MetaDto{}, err
	}
	return meme20Meta.ToDto(), nil
}

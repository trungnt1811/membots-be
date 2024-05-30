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

func (u *memeceptionUCase) GetMemeceptionBySymbol(ctx context.Context, symbol string) (dto.MemeceptionDetailResp, error) {
	memeMeta, err := u.MemeceptionRepository.GetMemeceptionBySymbol(ctx, symbol)
	if err != nil {
		return dto.MemeceptionDetailResp{}, err
	}
	// TODO: get price from RPC
	memeceptionDetailResp := dto.MemeceptionDetailResp{
		Meme:  memeMeta.ToDto(),
		Price: 0,
	}
	return memeceptionDetailResp, nil
}

func (u *memeceptionUCase) GetMemeceptions(ctx context.Context) (dto.MemeceptionResp, error) {
	// Get list meme past
	listMemePast, err := u.MemeceptionRepository.GetMemeceptionsPast(ctx)
	if err != nil {
		return dto.MemeceptionResp{}, err
	}
	listMemePastDto := make([]dto.MemeceptionCommon, 0)
	for _, meme := range listMemePast {
		listMemePastDto = append(listMemePastDto, meme.ToCommonRespDto())
	}

	// Get list meme upcoming
	listMemeUpcoming, err := u.MemeceptionRepository.GetMemeceptionsUpcoming(ctx)
	if err != nil {
		return dto.MemeceptionResp{}, err
	}
	listMemeUpcomingDto := make([]dto.MemeceptionCommon, 0)
	for _, meme := range listMemeUpcoming {
		listMemeUpcomingDto = append(listMemeUpcomingDto, meme.ToCommonRespDto())
	}

	// Get list meme live
	listMemeLive, err := u.MemeceptionRepository.GetMemeceptionsLive(ctx)
	if err != nil {
		return dto.MemeceptionResp{}, err
	}
	listMemeLiveDto := make([]dto.MemeceptionCommon, 0)
	for _, meme := range listMemeLive {
		listMemeLiveDto = append(listMemeLiveDto, meme.ToCommonRespDto())
	}

	memeceptionResp := dto.MemeceptionResp{
		Live:     listMemeLiveDto,
		Upcoming: listMemeUpcomingDto,
		Past:     listMemePastDto,
	}
	return memeceptionResp, nil
}

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

func (u *memeceptionUCase) GetMemeceptionByContractAddress(ctx context.Context, contractAddress string) (dto.MemeceptionDetailResp, error) {
	memeMeta, err := u.MemeceptionRepository.GetMemeceptionByContractAddress(ctx, contractAddress)
	if err != nil {
		return dto.MemeceptionDetailResp{}, err
	}
	// TODO: get price from RPC
	// TODO: get nfts info from graphnode in case meta is MEME404
	memeceptionDetailResp := dto.MemeceptionDetailResp{
		Meme:  memeMeta.ToDto(),
		Price: 0,
	}
	return memeceptionDetailResp, nil
}

func (u *memeceptionUCase) GetMemeceptions(ctx context.Context) (dto.MemeceptionsResp, error) {
	// Get list meme past
	listMemePast, err := u.MemeceptionRepository.GetMemeceptionsPast(ctx)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	listMemePastDto := make([]dto.MemeceptionCommon, 0)
	for _, meme := range listMemePast {
		listMemePastDto = append(listMemePastDto, meme.ToCommonDto())
	}

	// Get list meme upcoming
	listMemeUpcoming, err := u.MemeceptionRepository.GetMemeceptionsUpcoming(ctx)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	listMemeUpcomingDto := make([]dto.MemeceptionCommon, 0)
	for _, meme := range listMemeUpcoming {
		listMemeUpcomingDto = append(listMemeUpcomingDto, meme.ToCommonDto())
	}

	// Get list meme live
	listMemeLive, err := u.MemeceptionRepository.GetMemeceptionsLive(ctx)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	listMemeLiveDto := make([]dto.MemeceptionCommon, 0)
	for _, meme := range listMemeLive {
		listMemeLiveDto = append(listMemeLiveDto, meme.ToCommonDto())
	}

	memeceptionsResp := dto.MemeceptionsResp{
		Live:     listMemeLiveDto,
		Upcoming: listMemeUpcomingDto,
		Past:     listMemePastDto,
	}
	return memeceptionsResp, nil
}

package memeception

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
)

type memeceptionUCase struct {
	MemeceptionRepository interfaces.MemeceptionRepository
}

func NewMemeceptionUCase(memeceptionRepository interfaces.MemeceptionRepository) interfaces.MemeceptionUCase {
	return &memeceptionUCase{
		MemeceptionRepository: memeceptionRepository,
	}
}

func (u *memeceptionUCase) CreateMeme(ctx context.Context, payload dto.CreateMemePayload) error {
	swapFeeBps, err := strconv.ParseUint(payload.MemeInfo.SwapFeeBps, 0, 64)
	if err != nil {
		return fmt.Errorf("cannot convert SwapFeePct to uint64: %w", err)
	}

	vestingAllocBps, err := strconv.ParseUint(payload.MemeInfo.VestingAllocBps, 0, 64)
	if err != nil {
		return fmt.Errorf("cannot convert VestingAllocBps to uint64: %w", err)
	}

	networkID, err := strconv.ParseUint(payload.Blockchain.ChainID, 0, 64)
	if err != nil {
		return fmt.Errorf("cannot convert ChainID to uint64: %w", err)
	}

	startAt, err := time.Parse(time.RFC3339, payload.Memeception.StartAt)
	if err != nil {
		return fmt.Errorf("cannot parse StartAt to unix timestamp: %w", err)
	}

	social := model.Social{}
	for _, payload := range payload.Socials {
		social.Provider = payload.Provider
		social.Username = payload.Username
		social.PhotoURL = payload.PhotoURL
		social.URL = payload.URL
		social.DisplayName = payload.DisplayName
		break
	}

	memeModel := model.Meme{
		Name:            payload.MemeInfo.Name,
		Symbol:          payload.MemeInfo.Symbol,
		Description:     payload.MemeInfo.Description,
		LogoUrl:         payload.MemeInfo.LogoUrl,
		BannerUrl:       payload.MemeInfo.BannerUrl,
		CreatorAddress:  payload.Blockchain.CreatorAddress,
		SwapFeeBps:      swapFeeBps,
		VestingAllocBps: vestingAllocBps,
		Meta:            payload.MemeInfo.Meta,
		NetworkID:       networkID,
		Website:         payload.MemeInfo.Website,
		Memeception: model.Memeception{
			StartAt:   uint64(startAt.Unix()),
			Ama:       payload.MemeInfo.Ama,
			TargetETH: payload.MemeInfo.TargetETH,
			UpdatedAtEpoch: uint64(time.Now().Unix()),
		},
		Social: social,
	}
	err = u.MemeceptionRepository.CreateMeme(ctx, memeModel)
	if err != nil {
		return fmt.Errorf("cannot create meme: %w", err)
	}
	return nil
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
	listMemePastDto := make([]dto.Memeception, 0)
	for _, meme := range listMemePast {
		listMemePastDto = append(listMemePastDto, meme.ToDto())
	}

	// Get list meme live
	listMemeLive, err := u.MemeceptionRepository.GetMemeceptionsLive(ctx)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	listMemeLiveDto := make([]dto.Memeception, 0)
	for _, meme := range listMemeLive {
		listMemeLiveDto = append(listMemeLiveDto, meme.ToDto())
	}

	memeceptionsResp := dto.MemeceptionsResp{
		MemeceptionsByStatus: dto.MemeceptionByStatus{
			Past: listMemePastDto,
			Live: listMemeLiveDto,
		},
	}
	return memeceptionsResp, nil
}

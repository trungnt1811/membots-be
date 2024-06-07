package memeception

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
	"github.com/flexstack.ai/membots-be/internal/util"
)

type memeceptionUCase struct {
	MemeceptionRepository interfaces.MemeceptionRepository
	ClientMemeception     *subgraphclient.Client
}

func NewMemeceptionUCase(
	memeceptionRepository interfaces.MemeceptionRepository,
	clientMemeception *subgraphclient.Client,
) interfaces.MemeceptionUCase {
	return &memeceptionUCase{
		MemeceptionRepository: memeceptionRepository,
		ClientMemeception:     clientMemeception,
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
			StartAt:        uint64(startAt.Unix()),
			Ama:            payload.MemeInfo.Ama,
			TargetETH:      payload.MemeInfo.TargetETH,
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
	// TODO: get nfts info from graphnode in case meta is MEME404
	ethPrice, err := util.GetETHPrice()
	if err != nil {
		return dto.MemeceptionDetailResp{}, err
	}
	memeceptionDetailResp := dto.MemeceptionDetailResp{
		Meme:  memeMeta.ToDto(),
		Price: ethPrice,
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

	// Get list latest coins
	listLatest, err := u.MemeceptionRepository.GetMemeceptionsLatest(ctx)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	listLatestCoins := make([]dto.Memeception, 0)
	for _, meme := range listLatest {
		listLatestCoins = append(listLatestCoins, meme.ToDto())
	}

	ethPrice, err := util.GetETHPrice()
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}

	// Get top 5 latest launchpad txs
	requestOpts := &subgraphclient.RequestOptions{
		First: 5,
		IncludeFields: []string{
			"memeToken",
			"type",
			"user",
			"amountETH",
			"transactionHash",
		},
	}
	response, err := u.ClientMemeception.ListSwapHistories(ctx, requestOpts)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	launchpadTxs := make([]dto.LaunchpadTx, 0)
	var memeContractAddresses []string
	for _, meme := range response.MemecoinBuyExits {
		txType := "BUY"
		if meme.Type == "MemecoinExit" {
			txType = "SELL"
		}
		usdAmount, err := util.ConvertWeiToUSD(ethPrice, meme.AmountETH)
		if err != nil {
			return dto.MemeceptionsResp{}, err
		}
		memeContractAddresses = append(memeContractAddresses, meme.MemeToken)
		launchpadTxs = append(launchpadTxs, dto.LaunchpadTx{
			TxType:              txType,
			TxHash:              meme.TransactionHash,
			WalletAddress:       meme.User,
			MemeContractAddress: meme.MemeToken,
			AmountUSD:           usdAmount,
		})
	}

	// Update meme symbol and logo url
	mapMeme, err := u.MemeceptionRepository.GetMapMemeSymbolAndLogoURL(ctx, memeContractAddresses)
	if err != nil {
		return dto.MemeceptionsResp{}, err
	}
	for index, launchpadTx := range launchpadTxs {
		meme := mapMeme[launchpadTx.MemeContractAddress]
		launchpadTxs[index].Symbol = meme.Symbol
		launchpadTxs[index].LogoUrl = meme.LogoUrl
	}

	memeceptionsResp := dto.MemeceptionsResp{
		MemeceptionsByStatus: dto.MemeceptionByStatus{
			Past: listMemePastDto,
			Live: listMemeLiveDto,
		},
		LatestCoins:       listLatestCoins,
		LatestLaunchpadTx: launchpadTxs,
		Price:             ethPrice,
	}
	return memeceptionsResp, nil
}

func (u *memeceptionUCase) GetMemeceptionBySymbol(ctx context.Context, symbol string) (dto.MemeceptionDetailResp, error) {
	memeMeta, err := u.MemeceptionRepository.GetMemeceptionBySymbol(ctx, symbol)
	if err != nil {
		return dto.MemeceptionDetailResp{}, err
	}
	// TODO: get nfts info from graphnode in case meta is MEME404
	ethPrice, err := util.GetETHPrice()
	if err != nil {
		return dto.MemeceptionDetailResp{}, err
	}
	memeceptionDetailResp := dto.MemeceptionDetailResp{
		Meme:  memeMeta.ToDto(),
		Price: ethPrice,
	}
	return memeceptionDetailResp, nil
}

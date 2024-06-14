package launchpad

import (
	"context"
	"fmt"
	"strconv"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type launchpadUCase struct {
	Client                *subgraphclient.Client
	MemeceptionRepository interfaces.MemeceptionRepository
}

func NewLaunchpadUcase(client *subgraphclient.Client, memeceptionRepository interfaces.MemeceptionRepository) interfaces.LaunchpadUCase {
	return &launchpadUCase{Client: client, MemeceptionRepository: memeceptionRepository}
}

func (uc *launchpadUCase) GetHistory(ctx context.Context, address string) (dto.LaunchpadInfoResp, error) {
	requestOpts := &subgraphclient.RequestOptions{
		IncludeFields: []string{
			"*",
		},
	}

	response, err := uc.Client.GetSwapHistoryByMemeToken(ctx, address, requestOpts)
	if err != nil {
		return dto.LaunchpadInfoResp{}, err
	}
	// convert response to dto.SwapHistoryByAddressRsp
	var transactions []dto.Transaction
	for _, meme := range response.MemecoinBuyExits {
		timestamp, err := strconv.ParseUint(meme.BlockTimestamp, 10, 64)
		if err != nil {
			timestamp = 0
		}
		txType := "BUY"
		if meme.Type == "MemecoinExit" {
			txType = "SELL"
		}
		transactions = append(transactions, dto.Transaction{
			AmountETH:     meme.AmountETH,
			AmountMeme:    meme.AmountMeme,
			WalletAddress: meme.User,
			TxHash:        meme.TransactionHash,
			TxType:        txType,
			MemeID:        meme.ID,
			Epoch:         timestamp,
		})
	}

	if transactions == nil {
		transactions = []dto.Transaction{} // Assign an empty slice
	}

	lauchpadInfo := dto.LaunchpadInfoResp{LaunchpadInfo: dto.LaunchpadInfo{
		Transactions: transactions,
		Status:       "LIVE",
		TargetETH:    "0.1",
		CollectedETH: "0",
		TxCounter:    "0",
	}}

	memeInfo, err := uc.MemeceptionRepository.GetMemeceptionByContractAddress(ctx, address)
	if err != nil {
		return lauchpadInfo, nil
	}

	status := "LIVE"
	if memeInfo.Memeception.TargetETH <= memeInfo.Memeception.CollectedETH {
		status = "ENDED_SOLD_OUT"
	}

	lauchpadInfo.LaunchpadInfo.TargetETH = fmt.Sprintf("%f", memeInfo.Memeception.TargetETH)
	lauchpadInfo.LaunchpadInfo.CollectedETH = fmt.Sprintf("%f", memeInfo.Memeception.CollectedETH)
	lauchpadInfo.LaunchpadInfo.Status = status

	return lauchpadInfo, nil
}

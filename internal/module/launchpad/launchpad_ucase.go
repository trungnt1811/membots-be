package launchpad

import (
	"context"
	"strconv"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type launchpadUCase struct {
	Client *subgraphclient.Client
}

func NewLaunchpadUcase(client *subgraphclient.Client) interfaces.LaunchpadUCase {
	return &launchpadUCase{Client: client}
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

	return dto.LaunchpadInfoResp{LaunchpadInfo: dto.LaunchpadInfo{
		Transactions: transactions,
		Status:       "",
		TargetETH:    "",
		CollectedETH: "",
	}}, nil
}
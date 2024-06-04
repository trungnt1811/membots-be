package launchpad

import (
	"context"
	"strconv"

	unigraphclient "github.com/emersonmacro/go-uniswap-subgraph-client"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type launchpadUCase struct {
	Client *unigraphclient.Client
}

func NewLaunchpadUcase(client *unigraphclient.Client) interfaces.LaunchpadUCase {
	return &launchpadUCase{Client: client}
}

func (uc *launchpadUCase) GetHistory(ctx context.Context, address string) (dto.LaunchpadInfoResp, error) {
	requestOpts := &unigraphclient.RequestOptions{
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
	for _, meme := range response.MemecoinExits {
		timestamp, err := strconv.ParseUint(meme.BlockTimestamp, 10, 64)
		if err != nil {
			timestamp = 0
		}
		transactions = append(transactions, dto.Transaction{
			AmountETH:     meme.AmountETH,
			AmountMeme:    meme.AmountMeme,
			WalletAddress: meme.User,
			TxHash:        meme.TransactionHash,
			TxType:        "BUY",
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

package worker

import (
	"context"
	"strconv"

	"github.com/robfig/cron/v3"

	"github.com/flexstack.ai/membots-be/internal/constant"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
	"github.com/flexstack.ai/membots-be/internal/util/log"
)

type UpdateMemeOnchainWorker struct {
	Repo              interfaces.MemeceptionRepository
	ClientMemeception *subgraphclient.Client
	ClientSwap        *subgraphclient.Client
}

func NewUpdateMemeOnchainWorker(
	repo interfaces.MemeceptionRepository,
	clientMemeception *subgraphclient.Client,
	clientSwap *subgraphclient.Client,
) UpdateMemeOnchainWorker {
	return UpdateMemeOnchainWorker{
		Repo:              repo,
		ClientMemeception: clientMemeception,
		ClientSwap:        clientSwap,
	}
}

func (worker UpdateMemeOnchainWorker) RunJob() {
	// Create a new cron scheduler
	c := cron.New(cron.WithSeconds())

	// Add a job that runs every 15 seconds
	_, err := c.AddFunc("*/15 * * * * *", func() {
		worker.updateMemeOnchain(worker.Repo, worker.ClientMemeception, worker.ClientSwap)
	})
	if err != nil {
		log.LG.Infof("failed to run job updateMemeOnchain: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	// Keep the program running
	// select {}
}

func (worker UpdateMemeOnchainWorker) updateMemeOnchain(
	repo interfaces.MemeceptionRepository,
	clientMemeception *subgraphclient.Client,
	clientSwap *subgraphclient.Client,
) {
	ctx := context.Background()
	listMemeProcessing, err := repo.GetListMemeProcessing(ctx)
	if err != nil {
		log.LG.Infof("GetListMemeProcessing error: %v", err)
	}
	for _, memeProcessing := range listMemeProcessing {
		requestOpts := &subgraphclient.RequestOptions{
			First: 1,
			IncludeFields: []string{
				"id",
				"memeToken",
				"params_salt",
				"tiers.id",
				"tiers.nftId",
				"tiers.lowerId",
				"tiers.upperId",
				"tiers.nftSymbol",
				"tiers.nftName",
				"tiers.amountThreshold",
				"tiers.baseURL",
				"tiers.isFungible",
			},
		}
		response, err := clientMemeception.GetMemeCreatedsByCreatorAndSymbol(
			ctx, memeProcessing.CreatorAddress, memeProcessing.Symbol, requestOpts,
		)
		if err != nil {
			log.LG.Infof("GetMemeCreatedsByCreatorAndSymbol error: %v", err)
		}
		if len(response.MemeCreateds) == 0 {
			log.LG.Info("MemeCreateds len is 0")
			continue
		}
		// Get token's total supply and decimals
		requestOpts = &subgraphclient.RequestOptions{
			First: 1,
			IncludeFields: []string{
				"id",
				"totalSupply",
				"decimals",
			},
		}
		var tokenInfoResp *subgraphclient.ListTokensResponse
		tokenInfoResp, err = clientSwap.GetTokensByNameAndSymbol(
			ctx, memeProcessing.Name, memeProcessing.Symbol, requestOpts,
		)
		if err != nil {
			log.LG.Infof("GetTokensByNameAndSymbol error: %v", err)
		}
		if len(tokenInfoResp.Tokens) == 0 {
			log.LG.Info("Tokens len is 0")
			continue
		}
		// TODO: handle meme404 nft later
		decimals, err := strconv.ParseUint(tokenInfoResp.Tokens[0].Decimals, 10, 64)
		if err != nil {
			log.LG.Infof("Error parsing Decimals: %v", err)
		}
		meme := model.Meme{
			ID:              memeProcessing.ID,
			Salt:            response.MemeCreateds[0].Salt,
			ContractAddress: response.MemeCreateds[0].MemeToken,
			Status:          uint64(constant.SUCCEED),
			TotalSupply:     tokenInfoResp.Tokens[0].TotalSupply,
			Decimals:        decimals,
		}
		err = repo.UpdateMeme(ctx, meme)
		if err != nil {
			log.LG.Infof("UpdateMeme error: %v", err)
		}
	}
}

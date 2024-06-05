package worker

import (
	"context"

	"github.com/robfig/cron/v3"

	"github.com/flexstack.ai/membots-be/internal/constant"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
	"github.com/flexstack.ai/membots-be/internal/util/log"
)

type UpdateMemeOnchainWorker struct {
	Repo   interfaces.MemeceptionRepository
	Client *subgraphclient.Client
}

func NewUpdateMemeOnchainWorker(repo interfaces.MemeceptionRepository, client *subgraphclient.Client) UpdateMemeOnchainWorker {
	return UpdateMemeOnchainWorker{
		Repo:   repo,
		Client: client,
	}
}

func (worker UpdateMemeOnchainWorker) RunJob() {
	// Create a new cron scheduler
	c := cron.New(cron.WithSeconds())

	// Add a job that runs every 15 seconds
	_, err := c.AddFunc("*/15 * * * * *", func() {
		worker.updateMemeOnchain(worker.Repo, worker.Client)
	})
	if err != nil {
		log.LG.Infof("failed to run job updateMemeOnchain: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	// Keep the program running
	select {}
}

func (worker UpdateMemeOnchainWorker) updateMemeOnchain(repo interfaces.MemeceptionRepository, client *subgraphclient.Client) {
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
				"params_name",
				"params_symbol",
				"params_startAt",
				"params_salt",
				"params_creator",
				"params_targetETH",
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
		response, err := client.GetMemeCreatedsByCreatorAndSymbol(ctx, memeProcessing.CreatorAddress, memeProcessing.Symbol, requestOpts)
		if err != nil {
			log.LG.Infof("GetMemeCreatedsByCreatorAndSymbol error: %v", err)
		}
		// TODO: handle meme404 nft later
		meme := model.Meme{
			ID:              memeProcessing.ID,
			Salt:            response.MemeCreateds[0].Salt,
			ContractAddress: response.MemeCreateds[0].MemeToken,
			Status:          uint64(constant.SUCCEED),
		}
		err = repo.UpdateMeme(ctx, meme)
		if err != nil {
			log.LG.Infof("UpdateMeme error: %v", err)
		}
	}
}

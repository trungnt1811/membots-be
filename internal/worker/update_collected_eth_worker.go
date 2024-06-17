package worker

import (
	"context"
	"math/big"

	"github.com/robfig/cron/v3"

	"github.com/flexstack.ai/membots-be/internal/constant"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/model"
	"github.com/flexstack.ai/membots-be/internal/util"
	"github.com/flexstack.ai/membots-be/internal/util/log"
)

type UpdateCollectedETHWorker struct {
	Repo              interfaces.MemeceptionRepository
	MemeceptionClient *subgraphclient.Client
}

func NewUpdateCollectedETHWorker(
	repo interfaces.MemeceptionRepository,
	memeceptionClient *subgraphclient.Client,
) UpdateCollectedETHWorker {
	return UpdateCollectedETHWorker{
		Repo:              repo,
		MemeceptionClient: memeceptionClient,
	}
}

func (worker UpdateCollectedETHWorker) RunJob() {
	// Create a new cron scheduler
	c := cron.New(cron.WithSeconds())

	// Add a job that runs every 15 seconds
	_, err := c.AddFunc("*/15 * * * * *", func() {
		worker.updateCollectedETH(worker.Repo, worker.MemeceptionClient)
	})
	if err != nil {
		log.LG.Infof("failed to run job updateCollectedETH: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	// Block the current goroutine so that the cron job keeps running
	select {}
}

func (worker UpdateCollectedETHWorker) updateCollectedETH(
	repo interfaces.MemeceptionRepository,
	clientMemeception *subgraphclient.Client,
) {
	ctx := context.Background()
	listMemeLive, err := repo.GetListMemeLive(ctx)
	if err != nil {
		log.LG.Infof("GetListMemeLive error: %v", err)
	}
	for _, memeLive := range listMemeLive {
		requestOpts := &subgraphclient.RequestOptions{
			IncludeFields: []string{
				"*",
			},
		}
		response, err := clientMemeception.GetCollectedETHById(
			ctx, memeLive.ContractAddress, requestOpts,
		)
		if err != nil {
			log.LG.Infof("GetCollectedETHById: %s error: %v", memeLive.ContractAddress, err)
			continue
		}
		if len(response.CollectedETH.AmountETH) == 0 {
			log.LG.Infof("%s: CollectedETH is 0", memeLive.ContractAddress)
			continue
		}
		amountWei := new(big.Int)
		_, ok := amountWei.SetString(response.CollectedETH.AmountETH, 10)
		if !ok {
			log.LG.Infof("%s: Invalid CollectedETH: %v", memeLive.ContractAddress, response.CollectedETH.AmountETH)
			continue
		}
		amountETH := util.WeiToEther(amountWei)
		collectedETH, _ := amountETH.Float64()
		status := memeLive.Memeception.Status
		if collectedETH >= memeLive.Memeception.TargetETH {
			status = uint64(constant.ENDED_SOLD_OUT)
		}
		memeception := model.Memeception{
			ID:           memeLive.Memeception.ID,
			CollectedETH: collectedETH,
			Status:       status,
		}
		err = repo.UpdateMemeception(ctx, memeception)
		if err != nil {
			log.LG.Infof("UpdateMemeception error: %v", err)
		}
	}
}

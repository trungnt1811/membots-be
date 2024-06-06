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

	// Add a job that runs every 60 seconds
	_, err := c.AddFunc("*/60 * * * * *", func() {
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
			First: 1,
			IncludeFields: []string{
				"memeToken",
				"amountETH",
			},
		}
		response, err := clientMemeception.GetMemeLiquidityAddedsByContractAddress(
			ctx, memeLive.ContractAddress, requestOpts,
		)
		if err != nil {
			log.LG.Infof("GetMemeLiquidityAddedsByContractAddress error: %v", err)
			continue
		}
		if len(response.MemeLiquidityAddeds) == 0 {
			log.LG.Info("MemeLiquidityAddeds len is 0")
			continue
		}
		amountWei := new(big.Int)
		_, ok := amountWei.SetString(response.MemeLiquidityAddeds[0].AmountETH, 10)
		if !ok {
			log.LG.Infof("Invalid AmountETH: %v", response.MemeLiquidityAddeds[0].AmountETH)
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

package worker

import (
	"context"
	"math/big"
	"time"

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
	// Create a ticker that ticks every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		worker.updateCollectedETH(worker.Repo, worker.MemeceptionClient)
	}

	// Block the current goroutine so that the ticker keeps running
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
		amountWei := new(big.Int)
		_, ok := amountWei.SetString(response.CollectedETH.AmountETH, 10)
		if !ok {
			log.LG.Infof("%s: Invalid CollectedETH: %v", memeLive.ContractAddress, response.CollectedETH.AmountETH)
			continue
		}
		if amountWei.Int64() == 0 {
			log.LG.Infof("%s: CollectedETH is 0", memeLive.ContractAddress)
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

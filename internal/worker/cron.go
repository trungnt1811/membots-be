package worker

import (
	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/module/memeception"
)

func RegisterCronJobs(db *gorm.DB, memeceptionClient, swapClient *subgraphclient.Client) {
	memeRepo := memeception.NewMemeceptionRepository(db)

	// SECTION: Update meme onchain worker
	updateMemeOnchainWorker := NewUpdateMemeOnchainWorker(memeRepo, memeceptionClient, swapClient)
	go updateMemeOnchainWorker.RunJob()

	// SECTION: Update collected ETH worker
	updateCollectedETHWorker := NewUpdateCollectedETHWorker(memeRepo, memeceptionClient)
	go updateCollectedETHWorker.RunJob()
}

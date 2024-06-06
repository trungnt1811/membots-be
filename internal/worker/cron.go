package worker

import (
	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/module/memeception"
)

func RegisterCronJobs(db *gorm.DB, clientMemeception, clientSwap *subgraphclient.Client) {
	// SECTION: Update meme onchain worker
	memeRepo := memeception.NewMemeceptionRepository(db)
	updateMemeOnchainWorker := NewUpdateMemeOnchainWorker(memeRepo, clientMemeception, clientSwap)
	updateMemeOnchainWorker.RunJob()
}

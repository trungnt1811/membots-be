package worker

import (
	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/module/memeception"
)

func RegisterCronJobs(db *gorm.DB, client *subgraphclient.Client) {
	// SECTION: Update meme onchain worker
	memeRepo := memeception.NewMemeceptionRepository(db)
	updateMemeOnchainWorker := NewUpdateMemeOnchainWorker(memeRepo, client)
	updateMemeOnchainWorker.RunJob()
}

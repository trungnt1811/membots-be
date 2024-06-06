package worker

import (
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type UpdateCollectedETHWorker struct {
	Repo              interfaces.MemeceptionRepository
	MemeceptionClient *subgraphclient.Client
}

func NewUpdateCollectedETHWorker(
	repo interfaces.MemeceptionRepository,
	memeceptionClient *subgraphclient.Client,
) UpdateMemeOnchainWorker {
	return UpdateMemeOnchainWorker{
		Repo:              repo,
		MemeceptionClient: memeceptionClient,
	}
}

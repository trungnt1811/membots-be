package route

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/conf"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/module/launchpad"
	"github.com/flexstack.ai/membots-be/internal/module/memeception"
	"github.com/flexstack.ai/membots-be/internal/module/swap"
	"github.com/flexstack.ai/membots-be/internal/worker"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	// SECTION: Create redis client
	// redisClient := caching.NewCachingRepository(context.Background(), rdb)

	// SECTION: Create subgraph clients
	swapClient := subgraphclient.NewClient(
		"https://api.studio.thegraph.com/query/76502/membots-ai-v3-mvp/version/latest",
		nil,
	)
	memeClient := subgraphclient.NewClient(
		"https://api.studio.thegraph.com/query/76502/membots-ai-memeception-mvp/version/latest",
		nil,
	)

	appRouter := v1.Group("/truglymeme")

	// SECTION: memeception
	memeceptionRepository := memeception.NewMemeceptionRepository(db)
	// memeceptionCache := memeception.NewMemeceptionCacheRepository(memeceptionRepository, redisClient)
	memeceptionUCase := memeception.NewMemeceptionUCase(memeceptionRepository, memeClient)
	memeceptionHandler := memeception.NewMemeceptionHandler(memeceptionUCase)
	appRouter.GET("/meme", memeceptionHandler.GetMemeceptionByMemeAddress)
	appRouter.GET("/memeceptions", memeceptionHandler.GetMemeceptions)
	appRouter.POST("/memes", memeceptionHandler.CreateMeme)

	// SECTION: swaps
	swapUCase := swap.NewSwapUcase(swapClient)
	swapHandler := swap.NewSwapHandler(swapUCase)
	appRouter.GET("/swaps", swapHandler.GetSwapHistoryByAddress)
	appRouter.GET("/quote", swapHandler.GetSwapRouter)

	// SECTION: launchpad
	launchpadUCase := launchpad.NewLaunchpadUcase(memeClient)
	launchpadHandler := launchpad.NewLaunchpadHandler(launchpadUCase)
	appRouter.GET("/launchpad", launchpadHandler.GetHistoryByAddress)

	// SECTION: cronjob
	worker.RegisterCronJobs(db, memeClient, swapClient)
}

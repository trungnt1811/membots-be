package route

import (
	"context"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/conf"
	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/flexstack.ai/membots-be/internal/infra/subgraphclient"
	"github.com/flexstack.ai/membots-be/internal/module/launchpad"
	"github.com/flexstack.ai/membots-be/internal/module/memeception"
	"github.com/flexstack.ai/membots-be/internal/module/stats"
	"github.com/flexstack.ai/membots-be/internal/module/swap"
	"github.com/flexstack.ai/membots-be/internal/worker"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	// SECTION: Create redis client
	rdb := conf.RedisConn()
	redisClient := caching.NewCachingRepository(context.Background(), rdb)

	// SECTION: Create subgraph clients
	swapClient := subgraphclient.NewClient(
		config.Blockchain.SwapSubgraphURL,
		nil,
	)
	memeClient := subgraphclient.NewClient(
		config.Blockchain.MemeceptionSubgraphURL,
		nil,
	)

	appRouter := v1.Group("")

	// SECTION: memeception
	memeceptionRepository := memeception.NewMemeceptionRepository(db)
	memeceptionCache := memeception.NewMemeceptionCacheRepository(memeceptionRepository, redisClient)
	memeceptionUCase := memeception.NewMemeceptionUCase(memeceptionCache, memeClient)
	memeceptionHandler := memeception.NewMemeceptionHandler(memeceptionUCase)
	appRouter.GET("/meme", memeceptionHandler.GetMemeDetail)
	appRouter.GET("/memeceptions", memeceptionHandler.GetMemeceptions)
	appRouter.POST("/memes", memeceptionHandler.CreateMeme)

	// SECTION: swaps
	swapUCase := swap.NewSwapUcase(swapClient)
	swapCache := swap.NewSwapCacheUCase(swapUCase, redisClient)
	swapHandler := swap.NewSwapHandler(swapCache)
	appRouter.GET("/swaps", swapHandler.GetSwapHistoryByAddress)
	appRouter.GET("/quote", swapHandler.GetSwapRouter)

	// SECTION: launchpad
	launchpadUCase := launchpad.NewLaunchpadUcase(memeClient, memeceptionRepository)
	launchpadCache := launchpad.NewLaunchpadCacheUCase(launchpadUCase, redisClient)
	launchpadHandler := launchpad.NewLaunchpadHandler(launchpadCache)
	appRouter.GET("/launchpad", launchpadHandler.GetHistoryByAddress)

	// SECTION: stats
	statsUCase := stats.NewStatsUcase(swapClient, memeceptionRepository)
	statsCache := stats.NewStatsCacheUCase(statsUCase, redisClient)
	statsUCaseHandler := stats.NewStatsHandler(statsCache)
	appRouter.GET("/stats", statsUCaseHandler.GetStatsByMemeAddress)

	// SECTION: workers
	worker.RegisterWorkers(db, memeClient, swapClient, config.Blockchain.MemeceptionAddress)
}

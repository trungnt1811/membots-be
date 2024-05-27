package route

import (
	"context"

	"github.com/flexstack.ai/membots-be/conf"
	"github.com/flexstack.ai/membots-be/internal/app/fair_launch"
	"github.com/flexstack.ai/membots-be/internal/infra/caching"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	// SECTION: Create redis client
	rdb := conf.RedisConn()
	redisClient := caching.NewCachingRepository(context.Background(), rdb)

	// SECTION: fair launch module
	appRouter := v1.Group("/fair-launch")

	fairLaunchRepository := fair_launch.NewFairLauchRepository(db)
	fairLaunchCache := fair_launch.NewFairLauchCacheRepository(fairLaunchRepository, redisClient)
	fairLaunchUCase := fair_launch.NewFairLauchUCaseUCase(fairLaunchCache)

	fairLaunchHandler := fair_launch.NewFairLaunchHandler(fairLaunchUCase)
	appRouter.GET("/meme", fairLaunchHandler.GetMeme20MetaByTicker)
}

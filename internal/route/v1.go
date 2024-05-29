package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/flexstack.ai/membots-be/conf"
	"github.com/flexstack.ai/membots-be/internal/app/memeception"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	// SECTION: Create redis client

	// redisClient := caching.NewCachingRepository(context.Background(), rdb)

	// SECTION: memeception module
	appRouter := v1.Group("/memeception")

	memeceptionRepository := memeception.NewMemeceptionRepository(db)
	// memeceptionCache := memeception.NewMemeceptionCacheRepository(memeceptionRepository, redisClient)
	memeceptionUCase := memeception.NewMemeceptionUCase(memeceptionRepository)

	memeceptionHandler := memeception.NewMemeceptionHandler(memeceptionUCase)
	appRouter.GET("", memeceptionHandler.GetMemeceptionBySymbol)
}

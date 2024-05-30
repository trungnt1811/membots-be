package route

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/conf"
	"github.com/flexstack.ai/membots-be/internal/module/memeception"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	// SECTION: Create redis client

	// redisClient := caching.NewCachingRepository(context.Background(), rdb)

	// SECTION: truglymeme
	appRouter := v1.Group("/truglymeme")

	memeceptionRepository := memeception.NewMemeceptionRepository(db)
	// memeceptionCache := memeception.NewMemeceptionCacheRepository(memeceptionRepository, redisClient)
	memeceptionUCase := memeception.NewMemeceptionUCase(memeceptionRepository)

	memeceptionHandler := memeception.NewMemeceptionHandler(memeceptionUCase)
	appRouter.GET("/memeception", memeceptionHandler.GetMemeceptionBySymbol)
	appRouter.GET("/memeceptions", memeceptionHandler.GetMemeceptions)
}

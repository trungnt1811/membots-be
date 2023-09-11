package route

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	campaign3 "github.com/astraprotocol/affiliate-system/internal/app/campaign"
	campaign2 "github.com/astraprotocol/affiliate-system/internal/app/console/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/redeem"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/middleware"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB, channel *util.Channel) {
	v1 := r.Group("/api/v1")

	// SECTION: Create auth middleware
	jwtMiddleware := middleware.CreateJWTMiddleware(db)

	// SECTION: AT Module
	atRepo := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)

	// SECTION: Campaign and link
	campaignRepo := campaign3.NewCampaignRepository(db)
	campaignUsecase := campaign3.NewCampaignUsecase(campaignRepo, atRepo)
	campaignHandler := campaign3.NewCampaignHandler(campaignUsecase)
	campaignRoute := v1.Group("/campaign")
	campaignRoute.POST("/link", jwtMiddleware, campaignHandler.PostGenerateAffLink)

	// SECTION: Order Module and link

	// SECTION: Redeem module
	redeemRepo := redeem.NewRedeemRepository(db)
	redeemUsecase := redeem.NewRedeemUsecase(redeemRepo)
	redeemHandler := redeem.NewRedeemHandler(redeemUsecase)

	redeemRoute := v1.Group("/redeem")
	redeemRoute.POST("/request", jwtMiddleware, redeemHandler.PostRequestRedeem)

	consoleCampRepository := campaign2.NewConsoleCampaignRepository(db)
	consoleCampUCase := campaign2.NewCampaignUCase(consoleCampRepository)
	consoleCampHandler := campaign2.NewConsoleCampHandler(consoleCampUCase)
	consoleRouter := v1.Group("console")
	consoleRouter.GET("/aff-campaign", consoleCampHandler.GetAllCampaign)
	consoleRouter.PUT("/aff-campaign/:id", consoleCampHandler.UpdateCampaignInfo)

	// SECTION: Reward module
	rewardRepo := reward.NewRewardRepository(db)
	rewardUsecase := reward.NewRewardUsecase(rewardRepo)
	rewardHandler := reward.NewRewardHandler(rewardUsecase)

	rewardRouter := v1.Group("/rewards")
	rewardRouter.GET("/by-order-id", rewardHandler.GetRewardByOrderId)
	rewardRouter.GET("", rewardHandler.GetAllReward)
	rewardRouter.GET("/history", rewardHandler.GetRewardHistory)
}

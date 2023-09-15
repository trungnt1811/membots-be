package route

import (
	"context"
	bannerApp "github.com/astraprotocol/affiliate-system/internal/app/aff_banner_app"
	bannerConsole "github.com/astraprotocol/affiliate-system/internal/app/console/banner"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_camp_app"
	"github.com/astraprotocol/affiliate-system/internal/app/auth"
	campaign3 "github.com/astraprotocol/affiliate-system/internal/app/campaign"
	campaignConsole "github.com/astraprotocol/affiliate-system/internal/app/console/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/redeem"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/infra/shipping"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB, channel *util.Channel) {
	v1 := r.Group("/api/v1")

	// SECTION: Create redis client
	rdb := conf.RedisConn()
	redisClient := caching.NewCachingRepository(context.Background(), rdb)

	// SECTION: Create auth handler
	authHandler := auth.NewAuthUseCase(redisClient, config.CreatorAuthUrl, config.AppAuthUrl)

	// SECTION: AT Module
	atRepo := accesstrade.NewAccessTradeRepository(config.AccessTradeAPIKey, 3, 30)

	// SECTION: Shipping Reward Service
	shippingClientConf := shipping.ShippingClientConfig{
		BaseUrl: config.RewardShipping.BaseUrl,
		ApiKey:  config.RewardShipping.ApiKey,
	}
	shippingClient := shipping.NewShippingClient(shippingClientConf)

	// SECTION: Kafka Queue

	// SECTION: Campaign and link
	campaignRepo := campaign3.NewCampaignRepository(db)
	campaignUsecase := campaign3.NewCampaignUsecase(campaignRepo, atRepo)
	campaignHandler := campaign3.NewCampaignHandler(campaignUsecase)
	campaignRoute := v1.Group("/campaign")
	campaignRoute.POST("/link", authHandler.CheckUserHeader(), campaignHandler.PostGenerateAffLink)

	// SECTION: Order Module and link
	orderRepo := order.NewOrderRepository(db)
	orderUcase := order.NewOrderUcase(orderRepo, atRepo)
	orderHandler := order.NewOrderHandler(orderUcase)
	orderRoute := v1.Group("/order")
	orderRoute.POST("/post-back", orderHandler.PostBackOrderHandle)
	orderRoute.GET("/:id", orderHandler.GetOrderDetails)

	// SECTION: Redeem module
	redeemRepo := redeem.NewRedeemRepository(db)
	redeemUsecase := redeem.NewRedeemUsecase(redeemRepo)
	redeemHandler := redeem.NewRedeemHandler(redeemUsecase)

	redeemRoute := v1.Group("/redeem")
	redeemRoute.POST("/request", authHandler.CheckUserHeader(), redeemHandler.PostRequestRedeem)

	consoleRouter := v1.Group("console")

	consoleCampRepository := campaignConsole.NewConsoleCampaignRepository(db)
	consoleCampUCase := campaignConsole.NewCampaignUCase(consoleCampRepository)
	consoleCampHandler := campaignConsole.NewConsoleCampHandler(consoleCampUCase)
	consoleRouter.GET("/aff-campaign", consoleCampHandler.GetAllCampaign)
	consoleRouter.PUT("/aff-campaign/:id", authHandler.CheckAdminHeader(), consoleCampHandler.UpdateCampaignInfo)
	consoleRouter.GET("/aff-campaign/:id", consoleCampHandler.GetCampaignById)

	consoleBannerRepository := bannerConsole.NewConsoleBannerRepository(db)
	consoleBannerUCase := bannerConsole.NewBannerUCase(consoleBannerRepository, consoleCampRepository)
	consoleBannerHandler := bannerConsole.NewConsoleBannerHandler(consoleBannerUCase)

	consoleRouter.GET("/aff-banner", consoleBannerHandler.GetAllBanner)
	consoleRouter.PUT("/aff-banner/:id", authHandler.CheckAdminHeader(), consoleBannerHandler.UpdateBannerInfo)
	consoleRouter.GET("/aff-banner/:id", consoleBannerHandler.GetBannerById)
	consoleRouter.POST("/aff-banner", consoleBannerHandler.AddAffBanner)

	// SECTION: Reward module
	rewardRepo := reward.NewRewardRepository(db)
	rewardUsecase := reward.NewRewardUsecase(rewardRepo, orderRepo, shippingClient)
	rewardHandler := reward.NewRewardHandler(rewardUsecase)

	rewardRouter := v1.Group("/rewards")
	rewardRouter.GET("", rewardHandler.GetAllReward)
	rewardRouter.GET("/summary", rewardHandler.GetRewardSummary)
	rewardRouter.GET("/claims", rewardHandler.GetClaimHistory)
	rewardRouter.GET("/claims/:id", rewardHandler.GetClaimDetails)
	rewardRouter.POST("/claims", rewardHandler.ClaimReward)

	// SECTION: App module
	appRouter := v1.Group("/app")

	affCampAppRepository := aff_camp_app.NewAffCampAppRepository(db)
	affCampAppCache := aff_camp_app.NewAffCampAppCacheRepository(affCampAppRepository, redisClient)
	affCampAppService := aff_camp_app.NewAffCampAppService(affCampAppCache)
	affCampAppHandler := aff_camp_app.NewAffCampAppHandler(affCampAppService)
	appRouter.GET("/aff-campaign", affCampAppHandler.GetAllAffCampaign)
	appRouter.GET("/aff-campaign/:id", affCampAppHandler.GetAffCampaignById)

	affAppBannerRepo := bannerApp.NewAppBannerRepository(db)
	affAppBannerUCase := bannerApp.NewBannerUCase(affAppBannerRepo)
	affAppBannerHandler := bannerApp.NewAppBannerHandler(affAppBannerUCase)

	appRouter.GET("/aff-banner", affAppBannerHandler.GetAllBanner)
	appRouter.GET("/aff-banner/:id", affAppBannerHandler.GetBannerById)

}

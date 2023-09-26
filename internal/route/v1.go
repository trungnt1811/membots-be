package route

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/reward"

	bannerApp "github.com/astraprotocol/affiliate-system/internal/app/aff_banner_app"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_brand"
	categoryApp "github.com/astraprotocol/affiliate-system/internal/app/aff_category_app"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_search"
	bannerConsole "github.com/astraprotocol/affiliate-system/internal/app/console/banner"
	consoleOrder "github.com/astraprotocol/affiliate-system/internal/app/console/order"
	"github.com/astraprotocol/affiliate-system/internal/app/console/statistic"
	"github.com/astraprotocol/affiliate-system/internal/app/home_page"
	"github.com/astraprotocol/affiliate-system/internal/app/user_favorite_brand"
	"github.com/astraprotocol/affiliate-system/internal/middleware"
	"github.com/go-co-op/gocron"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_camp_app"
	"github.com/astraprotocol/affiliate-system/internal/app/auth"
	campaign1 "github.com/astraprotocol/affiliate-system/internal/app/campaign"
	campaignConsole "github.com/astraprotocol/affiliate-system/internal/app/console/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/redeem"
	"github.com/astraprotocol/affiliate-system/internal/app/user_view_aff_camp"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
	"github.com/astraprotocol/affiliate-system/internal/infra/shipping"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, config *conf.Configuration, db *gorm.DB) {
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
	userViewAffCampQueue := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_USER_VIEW_AFF_CAMP)

	// SECTION: Campaign and link
	campRepo := campaign1.NewCampaignRepository(db)
	campaignUCase := campaign1.NewCampaignUCase(campRepo, atRepo)
	campaignHandler := campaign1.NewCampaignHandler(campaignUCase)
	campaignRoute := v1.Group("/campaign")
	campaignRoute.POST("/link", authHandler.CheckUserHeader(), campaignHandler.PostGenerateAffLink)

	// SECTION: Order Module and link
	orderRepo := order.NewOrderRepository(db)
	orderUcase := order.NewOrderUCase(orderRepo, atRepo)
	orderHandler := order.NewOrderHandler(orderUcase)
	orderRoute := v1.Group("/order")
	orderRoute.POST("/post-back", orderHandler.PostBackOrderHandle)
	orderRoute.GET("", orderHandler.GetOrderHistory)
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

	consoleOrderRepo := consoleOrder.NewConsoleOrderRepository(db)
	consoleOrderUcase := consoleOrder.NewConsoleOrderUcase(consoleOrderRepo)
	consoleOrderHandler := consoleOrder.NewConsoleOrderHandler(consoleOrderUcase)

	// SECTION: Console Order
	consoleOrderRouter := consoleRouter.Group("orders", authHandler.CheckAdminHeader())
	consoleOrderRouter.GET("", consoleOrderHandler.GetOrderList)
	consoleOrderRouter.GET("/:orderId", consoleOrderHandler.GetOrderByOrderId)

	// SECTION: Console Summary
	statisticRepo := statistic.NewStatisticRepository(db)
	statisticUcase := statistic.NewStatisticUCase(statisticRepo)
	statisticHandler := statistic.NewStatisticHandler(statisticUcase)
	consoleRouter.GET("/summary", authHandler.CheckAdminHeader(), statisticHandler.GetSummary)
	consoleRouter.GET("/summary/:campaignId", authHandler.CheckAdminHeader(), statisticHandler.GetCampaignSummary)

	// SECTION: App module
	appRouter := v1.Group("/app")

	streamChannel := make(chan []*dto.UserViewAffCampDto, 1024)
	userViewAffCampProducer := msgqueue.NewUserViewAffCampProducer(userViewAffCampQueue, streamChannel)
	userViewAffCampProducer.Start()

	userViewAffCampRepository := user_view_aff_camp.NewUserViewAffCampRepository(db)
	userViewAffCampCache := user_view_aff_camp.NewUserViewAffCampCacheRepository(userViewAffCampRepository, redisClient)
	userViewAffCampUCase := user_view_aff_camp.NewUserViewAffCampUCase(userViewAffCampCache)

	userFavoriteBrandRepository := user_favorite_brand.NewUserFavoriteBrandRepository(db)
	userFavoriteBrandCache := user_favorite_brand.NewUserFavoriteBrandCacheRepository(userFavoriteBrandRepository, redisClient)

	affCampAppRepository := aff_camp_app.NewAffCampAppRepository(db)
	affCampAppCache := aff_camp_app.NewAffCampAppCacheRepository(affCampAppRepository, redisClient)
	affBrandRepository := aff_brand.NewAffBrandRepository(db)
	affBrandCache := aff_brand.NewAffBrandCacheRepository(affBrandRepository, redisClient)

	affCampAppUCase := aff_camp_app.NewAffCampAppUCase(affCampAppCache, affBrandCache, userFavoriteBrandCache, streamChannel)
	affCampAppHandler := aff_camp_app.NewAffCampAppHandler(affCampAppUCase)
	appRouter.GET("/aff-campaign", authHandler.CheckUserHeader(), affCampAppHandler.GetAllAffCampaign)
	appRouter.GET("/aff-campaign/:id", authHandler.CheckUserHeader(), affCampAppHandler.GetAffCampaignById)

	affBrandUCase := aff_brand.NewAffBrandUCase(affBrandCache, affCampAppCache, userFavoriteBrandCache)
	affBrandHandler := aff_brand.NewAffBrandHandler(userViewAffCampUCase, affBrandUCase)
	appRouter.GET("brand", authHandler.CheckUserHeader(), affBrandHandler.GetListAffBrandByUser)

	homePageUCase := home_page.NewHomePageUCase(affBrandCache, affCampAppCache, userFavoriteBrandCache, userViewAffCampCache)
	homePageHandler := home_page.NewHomePageHandler(homePageUCase)
	appRouter.GET("/home-page", authHandler.CheckUserHeader(), homePageHandler.GetHomePage)

	affAppBannerRepo := bannerApp.NewAppBannerRepository(db)
	affAppBannerUCase := bannerApp.NewBannerUCase(affAppBannerRepo)
	affAppBannerHandler := bannerApp.NewAppBannerHandler(affAppBannerUCase)

	appRouter.GET("/aff-banner", affAppBannerHandler.GetAllBanner)
	appRouter.GET("/aff-banner/:id", affAppBannerHandler.GetBannerById)

	affCategoryRepo := categoryApp.NewAppCategoryRepository(db)
	affCategoryUCase := categoryApp.NewAffCategoryUCase(affCategoryRepo)
	affCategoryHandler := categoryApp.NewAffCategoryHandler(affCategoryUCase)
	appRouter.GET("/aff-categories", affCategoryHandler.GetAllCategory)
	appRouter.GET("/aff-categories/:categoryId", affCategoryHandler.GetAllAffCampaignInCategory)

	affSearchRepo := aff_search.NewAffSearchRepository(db)
	affSearchUCase := aff_search.NewAffSearchUCase(affSearchRepo)
	affSearchHandler := aff_search.NewAffSearchHandler(affSearchUCase)
	appRouter.GET("/aff-search", affSearchHandler.AffSearch)

	// SECTION: Reward module
	rewardConf := reward.RewardConfig{
		SellerId:      config.Aff.SellerId,
		RewardProgram: config.Aff.RewardProgram,
	}
	rewardRepo := reward.NewRewardRepository(db)
	rewardUsecase := reward.NewRewardUCase(rewardRepo, orderRepo, shippingClient, rewardConf)
	rewardHandler := reward.NewRewardHandler(rewardUsecase)
	withdrawRateLimit := middleware.NewRateLimit(rdb, 3*time.Second, 1)

	rewardRouter := appRouter.Group("/rewards", authHandler.CheckUserHeader())
	rewardRouter.GET("/summary", rewardHandler.GetRewardSummary)
	rewardRouter.GET("/withdraw", rewardHandler.GetWithdrawHistory)
	rewardRouter.POST("/withdraw", withdrawRateLimit, rewardHandler.WithdrawReward)

	orderRouteApp := appRouter.Group("/orders", authHandler.CheckUserHeader())
	orderRouteApp.GET("", orderHandler.GetOrderHistory)
	orderRouteApp.GET("/:id", orderHandler.GetOrderDetails)

	// SECTION: Cron jobs
	cron := gocron.NewScheduler(time.UTC)
	_, err := cron.Every(5).Minute().Do(func() {
		affBrandUCase.UpdateCacheListCountFavouriteAffBrand(context.Background())
	})
	if err == nil {
		cron.StartAsync()
	}
}

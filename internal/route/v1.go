package route

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/aff_brand"

	bannerApp "github.com/astraprotocol/affiliate-system/internal/app/aff_banner_app"
	categoryApp "github.com/astraprotocol/affiliate-system/internal/app/aff_category_app"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_search"
	bannerConsole "github.com/astraprotocol/affiliate-system/internal/app/console/banner"
	consoleOrder "github.com/astraprotocol/affiliate-system/internal/app/console/order"
	"github.com/go-co-op/gocron"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_camp_app"
	"github.com/astraprotocol/affiliate-system/internal/app/auth"
	campaign3 "github.com/astraprotocol/affiliate-system/internal/app/campaign"
	campaignConsole "github.com/astraprotocol/affiliate-system/internal/app/console/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/redeem"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/app/user_view_aff_camp"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
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
	userViewAffCampQueue := msgqueue.NewKafkaProducer(msgqueue.KAFKA_TOPIC_USER_VIEW_AFF_CAMP)

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
	consoleOrderUcase := consoleOrder.NewOrderUcase(consoleOrderRepo)
	consoleOrderHandler := consoleOrder.NewConsoleOrderHandler(consoleOrderUcase)

	// SECTION: Console Order
	consoleOrderRouter := consoleRouter.Group("orders", authHandler.CheckAdminHeader())
	consoleOrderRouter.GET("", consoleOrderHandler.GetOrderList)

	// SECTION: Reward module
	rewardRepo := reward.NewRewardRepository(db)
	rewardUsecase := reward.NewRewardUsecase(rewardRepo, orderRepo, shippingClient)
	rewardHandler := reward.NewRewardHandler(rewardUsecase)

	rewardRouter := v1.Group("/rewards")
	rewardRouter.GET("/summary", rewardHandler.GetRewardSummary)
	rewardRouter.GET("/withdraw", rewardHandler.GetWithdrawHistory)
	rewardRouter.GET("/withdraw/:id", rewardHandler.GetWithdrawDetails)
	rewardRouter.POST("/withdraw", rewardHandler.WithdrawReward)

	// SECTION: App module
	streamChannel := make(chan []*dto.UserViewAffCampDto, 1024)

	appRouter := v1.Group("/app")

	affCampAppRepository := aff_camp_app.NewAffCampAppRepository(db)
	affCampAppCache := aff_camp_app.NewAffCampAppCacheRepository(affCampAppRepository, redisClient)
	affCampAppUCase := aff_camp_app.NewAffCampAppUCase(affCampAppCache, streamChannel)
	affCampAppHandler := aff_camp_app.NewAffCampAppHandler(affCampAppUCase)
	appRouter.GET("/aff-campaign", authHandler.CheckUserHeader(), affCampAppHandler.GetAllAffCampaign)
	appRouter.GET("/aff-campaign/:id", authHandler.CheckUserHeader(), affCampAppHandler.GetAffCampaignById)

	userViewAffCampProducer := msgqueue.NewUserViewAffCampProducer(userViewAffCampQueue, streamChannel)
	userViewAffCampProducer.Start()

	userViewAffCampRepository := user_view_aff_camp.NewUserViewAffCampRepository(db)
	userViewAffCampCache := user_view_aff_camp.NewUserViewAffCampCacheRepository(userViewAffCampRepository, redisClient)
	userViewAffCampUCase := user_view_aff_camp.NewUserViewAffCampUCase(userViewAffCampCache)
	userViewAffCampHandler := user_view_aff_camp.NewUserViewAffCampHandler(userViewAffCampUCase)
	appRouter.GET("brand/recently-visited-section", authHandler.CheckUserHeader(), userViewAffCampHandler.GetListRecentlyVisitedSection)

	affBrandRepository := aff_brand.NewAffBrandRepository(db)
	affBrandCache := aff_brand.NewAffBrandCacheRepository(affBrandRepository, redisClient)
	affBrandUCase := aff_brand.NewAffBrandUCase(affBrandCache, affCampAppCache)
	affBrandHandler := aff_brand.NewAffBrandHandler(affBrandUCase)
	appRouter.GET("brand/top-favorited", affBrandHandler.GetTopFavouriteAffBrand)

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
	consoleRouter.GET("/aff-search", authHandler.CheckAdminHeader(), affSearchHandler.SearchConsole)

	// SECTION: Cron jobs
	cron := gocron.NewScheduler(time.UTC)
	_, _ = cron.Every(5).Minute().Do(func() {
		affBrandUCase.UpdateCacheListCountFavouriteAffBrand(context.Background())
	})
	cron.StartAsync()
}

package route

import (
	"context"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/aff_camp_app"
	"github.com/astraprotocol/affiliate-system/internal/app/auth"
	campaign3 "github.com/astraprotocol/affiliate-system/internal/app/campaign"
	campaign2 "github.com/astraprotocol/affiliate-system/internal/app/console/campaign"
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

	// SECTION: Redeem module
	redeemRepo := redeem.NewRedeemRepository(db)
	redeemUsecase := redeem.NewRedeemUsecase(redeemRepo)
	redeemHandler := redeem.NewRedeemHandler(redeemUsecase)

	redeemRoute := v1.Group("/redeem")
	redeemRoute.POST("/request", authHandler.CheckUserHeader(), redeemHandler.PostRequestRedeem)

	consoleCampRepository := campaign2.NewConsoleCampaignRepository(db)
	consoleCampUCase := campaign2.NewCampaignUCase(consoleCampRepository)
	consoleCampHandler := campaign2.NewConsoleCampHandler(consoleCampUCase)
	consoleRouter := v1.Group("console")
	consoleRouter.GET("/aff-campaign", consoleCampHandler.GetAllCampaign)
	consoleRouter.PUT("/aff-campaign/:id", consoleCampHandler.UpdateCampaignInfo)
	consoleRouter.GET("/aff-campaign/:id", consoleCampHandler.GetCampaignById)

	// SECTION: Reward module
	rewardRepo := reward.NewRewardRepository(db)
	rewardUsecase := reward.NewRewardUsecase(rewardRepo, orderRepo, shippingClient)
	rewardHandler := reward.NewRewardHandler(rewardUsecase)

	rewardRouter := v1.Group("/rewards")
	rewardRouter.GET("/by-order-id", rewardHandler.GetRewardByOrderId)
	rewardRouter.GET("", rewardHandler.GetAllReward)
	rewardRouter.GET("/summary", rewardHandler.GetRewardSummary)
	rewardRouter.GET("/claims", rewardHandler.GetClaimHistory)
	rewardRouter.POST("/claims", rewardHandler.ClaimReward)

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
	userViewAffCampUCase := user_view_aff_camp.NewUserViewAffCampUCase(userViewAffCampRepository)
	userViewAffCampHandler := user_view_aff_camp.NewUserViewAffCampHandler(userViewAffCampUCase)
	appRouter.GET("/recently-visited-section", authHandler.CheckUserHeader(), userViewAffCampHandler.GetListRecentlyVisitedSection)
}

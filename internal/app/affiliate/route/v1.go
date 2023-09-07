package route

import (
	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/middleware"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/redeem"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util"
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
	campaignRepo := campaign.NewCampaignRepository(db)
	campaignUsecase := campaign.NewCampaignUsecase(campaignRepo, atRepo)
	campaignHandler := campaign.NewCampaignHandler(campaignUsecase)
	campaignRoute := v1.Group("/campaign")
	campaignRoute.POST("/link", jwtMiddleware, campaignHandler.PostGenerateAffLink)

	// SECTION: Order Module and link

	// SECTION: Redeem module
	redeemRepo := redeem.NewRedeemRepository(db)
	redeemUsecase := redeem.NewRedeemUsecase(redeemRepo)
	redeemHandler := redeem.NewRedeemHandler(redeemUsecase)

	redeemRoute := v1.Group("/redeem")
	redeemRoute.POST("/request", jwtMiddleware, redeemHandler.PostRequestRedeem)
}

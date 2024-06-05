package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/gorm/logger"

	pagination "github.com/AstraProtocol/reward-libs/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/flexstack.ai/membots-be/conf"
	"github.com/flexstack.ai/membots-be/internal/middleware"
	routeV1 "github.com/flexstack.ai/membots-be/internal/route"
	"github.com/flexstack.ai/membots-be/internal/util/log"
)

// @title           membots-be API
// @version         1.0
// @description     This Swagger docs for membots-be.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

func RunApp(config *conf.Configuration) {
	// Set zerolog global level
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.LG = log.NewZerologLogger(os.Stdout, zerolog.InfoLevel)

	if config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(pagination.Default())

	db := conf.DBConnWithLoglevel(logger.Info)

	// SECTION: Register middlewares
	// r.Use(middleware.StructuredLogger())
	r.Use(middleware.RequestLogger(log.LG.Instance))
	r.Use(gin.Recovery())

	// SECTION: Register routes
	routeV1.RegisterRoutes(r, config, db)

	// SECTION: Register general handlers
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%s is still alive", config.AppName),
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// SECTION: Run Gin router
	err := r.Run(fmt.Sprintf("0.0.0.0:%v", config.AppPort))
	if err != nil {
		log.LG.Fatalf("failed to run gin router: %v", err)
	}

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}

package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/astraprotocol/affiliate-system/internal/middleware"
	routeV1 "github.com/astraprotocol/affiliate-system/internal/route"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/astraprotocol/affiliate-system/internal/webhook"
	"github.com/astraprotocol/affiliate-system/internal/worker/cron"
	"github.com/astraprotocol/affiliate-system/internal/worker/kafkaconsumer"

	pagination "github.com/AstraProtocol/reward-libs/middleware"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/gin-gonic/gin"
	ginPrometheus "github.com/mcuadros/go-gin-prometheus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Affiliate System API
// @version         1.0
// @description     This Swagger docs for Astra Affiliate System.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description			Use for authorization of reward creator

// RunApp @securityDefinitions.apiKey	BasicKeyAuth
// @in							header
// @name						Authorization
// @description			Use for authorization during server to server calls
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

	// SECTION : channels
	channel := util.NewChannel()

	// SECTION: Register middlewares
	// r.Use(middleware.StructuredLogger())
	r.Use(middleware.RequestLogger(log.LG.Instance))
	r.Use(gin.Recovery())

	// SECTION: Register routes
	routeV1.RegisterRoutes(r, config, db, channel)

	// SECTION: Register general handlers
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%s is still alive", config.AppName),
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	p := ginPrometheus.NewPrometheus(config.AppName)
	p.Use(r)

	var err error
	// SECTION: Run Discord Webhook
	webhook.Whm, err = webhook.NewWebHookManagerFromConfig(config.Webhook)
	if err != nil {
		log.LG.Fatalf("failed to init webhook %v", err)
	}
	err = webhook.Whm.Start()
	if err != nil {
		log.LG.Fatalf("failed to start webhook %v", err)
	}

	// SECTION: Run worker
	cron.RegisterCronJobs(config, db)

	// SECTION: Run kafka consumer
	kafkaconsumer.RegisConsumers(config, db)

	// SECTION: Run Gin router
	err = r.Run(fmt.Sprintf("0.0.0.0:%v", config.AppPort))
	if err != nil {
		log.LG.Fatalf("failed to run gin router %v", err)
	}

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}

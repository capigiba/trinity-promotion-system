package router

import (
	_ "trinity/docs"
	"trinity/internal/initialize"
	"trinity/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var log = logger.NewLogger("router")

func SetupRouter(app *initialize.App) *gin.Engine {
	router := gin.Default()

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register routes
	api := router.Group("/")

	// Campaign routes
	campaignRoutes := api.Group("/campaigns")
	app.CampaignHandler.RegisterRoutes(campaignRoutes)

	// Voucher routes
	voucherRoutes := api.Group("/vouchers")
	app.VoucherHandler.RegisterRoutes(voucherRoutes)

	// Purchase routes
	purchaseRoutes := api.Group("/purchases")
	app.PurchaseHandler.RegisterRoutes(purchaseRoutes)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	log.Info("Routes have been set up")

	return router
}

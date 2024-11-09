package initialize

import (
	"trinity/config"
	"trinity/internal/campaign"
	"trinity/internal/database"
	"trinity/internal/purchase"
	"trinity/internal/subscription"
	"trinity/internal/voucher"
	"trinity/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

var log logger.Logger

type App struct {
	DB              *mongo.Database
	CampaignHandler *campaign.Handler
	VoucherHandler  *voucher.Handler
	PurchaseHandler *purchase.Handler
}

// Initialize sets up the application dependencies
func Initialize(cfg *config.Config) (*App, error) {
	log = logger.NewLogger("Initializer")

	db, err := database.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Errorf("failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Setup indexes
	err = database.SetupIndexes(db)
	if err != nil {
		log.Errorf("failed to set up indexes: %v", err)
		return nil, err
	}

	// Repositories
	campaignRepo := campaign.NewRepository(db)
	voucherRepo := voucher.NewRepository(db)
	subscriptionRepo := subscription.NewRepository(db)
	purchaseRepo := purchase.NewRepository(db)

	// Services
	campaignService := campaign.NewService(campaignRepo, voucherRepo)
	voucherService := voucher.NewService(voucherRepo)
	purchaseService := purchase.NewService(purchaseRepo, voucherRepo, subscriptionRepo)

	// Handlers
	campaignHandler := campaign.NewHandler(campaignService)
	voucherHandler := voucher.NewHandler(voucherService)
	purchaseHandler := purchase.NewHandler(purchaseService)

	app := &App{
		DB:              db,
		CampaignHandler: campaignHandler,
		VoucherHandler:  voucherHandler,
		PurchaseHandler: purchaseHandler,
	}

	return app, nil
}

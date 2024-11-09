package main

import (
	"trinity/config"
	"trinity/internal/initialize"
	"trinity/internal/router"
	"trinity/pkg/localization"
	"trinity/pkg/logger"

	"github.com/joho/godotenv"
)

var log logger.Logger

// @title Trinity App API
// @version 1.0
// @description API documentation for Trinity App

// @contact.name API Support
// @contact.url http://www.capigiba.com/support
// @contact.email giabao.vonguyen@gmail.com

// @host localhost:8080
// @BasePath /
func main() {
	log = logger.NewLogger("main")
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found")
	}

	// Load configuration
	cfg := config.LoadConfig()
	log.Info("Configuration loaded")

	// Initialize localization after config is loaded
	localization.Initialize()

	// Initialize application dependencies
	app, err := initialize.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	log.Info("Application initialized")

	// Set up the router
	r := router.SetupRouter(app)
	log.Info("Router set up")

	// Start the server
	log.Infof("Server is running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

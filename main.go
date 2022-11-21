package main

import (
	"spotit-backend/config"
	"spotit-backend/infra/database"
	"spotit-backend/infra/logger"
	"spotit-backend/routers"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// Set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Europe/Helsinki")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	dbDSN := config.DbConfiguration()

	if err := database.DbConnection(dbDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.RunTLS(config.ServerConfig(), "/etc/letsencrypt/live/spotit.codebite.fi/fullchain.pem", "/etc/letsencrypt/live/spotit.codebite.fi/privkey.pem"))
	// logger.Fatalf("%v", router.Run(config.ServerConfig()))
}

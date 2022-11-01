package main

import (
	"gin-boilerplate/config"
	"gin-boilerplate/infra/database"
	"gin-boilerplate/infra/logger"
	"gin-boilerplate/routers"
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
	logger.Fatalf("%v", router.Run(config.ServerConfig()))
}

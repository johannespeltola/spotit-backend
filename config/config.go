package config

import (
	"fmt"
	"spotit-backend/infra/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file, %s", err)
		return err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return err
	}

	return nil
}

func GetEntsoeBase() string {
	return fmt.Sprintf("https://transparency.entsoe.eu/api?securityToken=%v&documentType=%v&in_Domain=%v&out_Domain=%v", viper.GetString("ENTSOE_TOKEN"), "A44", viper.GetString("ENTSOE_DOMAIN"), viper.GetString("ENTSOE_DOMAIN"))
}

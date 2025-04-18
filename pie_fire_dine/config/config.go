package config

import (
	"os"

	"pie_fire_dine/logger"

	"github.com/spf13/viper"
)

var config Config

type Config struct {
	appName              string
	port                 string
	log                  *logger.LogConfig
	csvMeatCategoryPath  string
	defaultSourceTextURL string
}

func Load() {
	viper.SetDefault("APP_PORT", "8080")

	viper.SetConfigName("config")
	if os.Getenv("ENVIRONMENT") == "test" {
		viper.SetConfigName("test")
	}

	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")

	// Read config from environment variables
	viper.AutomaticEnv()
	// Read config from config file if exists (override environment variables)
	_ = viper.ReadInConfig()

	config = Config{
		appName: extractStringValue("APP_NAME"),
		port:    extractStringValue("APP_PORT"),
		log:     newLoggerConfig(),

		csvMeatCategoryPath:  extractStringValue("CSV_MEAT_CATEGORY_PATH"),
		defaultSourceTextURL: extractStringValue("DEFAULT_SOURCE_TEXT_URL"),
	}
}

func newLoggerConfig() *logger.LogConfig {
	return &logger.LogConfig{
		Out:    os.Stdout,
		Level:  extractStringValue("LOG_LEVEL"),
		Format: extractStringValue("LOG_FORMAT"),
	}
}

func GetAppName() string {
	return config.appName
}

func GetPort() string {
	return config.port
}

func GetLogger() *logger.LogConfig {
	return config.log
}

func GetCsvMeatCategoryPath() string {
	return config.csvMeatCategoryPath
}

func GetDefaultSourceTextURL() string {
	return config.defaultSourceTextURL
}

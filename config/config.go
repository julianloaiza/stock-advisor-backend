package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Address           string
	DatabaseURL       string
	StockAPIURL       string
	StockAPIKey       string
	SyncMaxIterations int
	SyncTimeout       int
}

func New() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	// Intenta cargar el archivo .env; si falla, se seguirán usando las variables de entorno
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ No se pudo leer .env, se utilizarán las variables de entorno")
	}

	viper.SetDefault("ADDRESS", ":8080")
	viper.SetDefault("SYNC_MAX_ITERATIONS", 100)

	config := &Config{
		Address:           viper.GetString("ADDRESS"),
		DatabaseURL:       viper.GetString("DATABASE_URL"),
		StockAPIURL:       viper.GetString("STOCK_API_URL"),
		StockAPIKey:       viper.GetString("STOCK_API_KEY"),
		SyncMaxIterations: viper.GetInt("SYNC_MAX_ITERATIONS"),
		SyncTimeout:       viper.GetInt("SYNC_TIMEOUT"),
	}

	return config
}

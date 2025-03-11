package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Address     string
	DatabaseURL string
	StockAPIURL string
	StockAPIKey string
	SimulateDB  bool
}

func New() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	// Intenta cargar el archivo .env; si falla, se seguirán usando las variables de entorno
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ No se pudo leer .env, se utilizarán las variables de entorno")
	}

	// Definir defaults si es necesario
	viper.SetDefault("ADDRESS", ":8080")
	viper.SetDefault("SIMULATE_DB", false)

	config := &Config{
		Address:     viper.GetString("ADDRESS"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		StockAPIURL: viper.GetString("STOCK_API_URL"),
		StockAPIKey: viper.GetString("STOCK_API_KEY"),
		SimulateDB:  viper.GetBool("SIMULATE_DB"),
	}

	// Si no se simula la BD, validar que DATABASE_URL esté definido
	if !config.SimulateDB {
		if config.DatabaseURL == "" {
			log.Fatal("❌ DATABASE_URL no está definido")
		}
		if config.StockAPIURL == "" || config.StockAPIKey == "" {
			log.Fatal("❌ Configuración de StockAPI incompleta")
		}
	}

	return config
}

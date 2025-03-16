package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config contiene la configuración de la aplicación.
type Config struct {
	Address            string
	DatabaseURL        string
	StockAPIURL        string
	StockAPIKey        string
	SyncMaxIterations  int
	SyncTimeout        int
	CORSAllowedOrigins string
}

// New crea una nueva instancia de Config.
func New() *Config {
	// Configuración de Viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Intentar cargar archivo .env
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ No se pudo leer .env, se utilizarán las variables de entorno")
	} else {
		log.Println("✅ Archivo .env cargado correctamente")
	}

	// Valores por defecto
	viper.SetDefault("ADDRESS", ":8080")
	viper.SetDefault("SYNC_MAX_ITERATIONS", 100)
	viper.SetDefault("SYNC_TIMEOUT", 60)
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")

	// Crear configuración
	config := &Config{
		Address:            viper.GetString("ADDRESS"),
		DatabaseURL:        viper.GetString("DATABASE_URL"),
		StockAPIURL:        viper.GetString("STOCK_API_URL"),
		StockAPIKey:        viper.GetString("STOCK_API_KEY"),
		SyncMaxIterations:  viper.GetInt("SYNC_MAX_ITERATIONS"),
		SyncTimeout:        viper.GetInt("SYNC_TIMEOUT"),
		CORSAllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
	}

	// Validar configuración
	if err := validateConfig(config); err != nil {
		log.Fatalf("❌ Error en la configuración: %v", err)
	}

	// Mostrar configuración (sin datos sensibles)
	logConfig(config)

	return config
}

// validateConfig verifica que los valores críticos no estén vacíos.
func validateConfig(cfg *Config) error {
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL no puede estar vacío")
	}
	if cfg.StockAPIURL == "" {
		return errors.New("STOCK_API_URL no puede estar vacío")
	}
	if cfg.StockAPIKey == "" {
		return errors.New("STOCK_API_KEY no puede estar vacío")
	}
	if cfg.SyncTimeout <= 0 {
		return errors.New("SYNC_TIMEOUT debe ser mayor que 0")
	}
	return nil
}

// logConfig muestra la configuración actual (sin datos sensibles).
func logConfig(cfg *Config) {
	log.Println("📋 Configuración cargada:")
	log.Printf("   - Servidor: %s", cfg.Address)
	log.Printf("   - DB: %s", maskString(cfg.DatabaseURL))
	log.Printf("   - API URL: %s", cfg.StockAPIURL)
	log.Printf("   - Max Iteraciones: %d", cfg.SyncMaxIterations)
	log.Printf("   - Timeout: %d segundos", cfg.SyncTimeout)
	log.Printf("   - CORS: %s", cfg.CORSAllowedOrigins)
}

// maskString oculta parte de una cadena para seguridad.
func maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return fmt.Sprintf("%s***", s[:8])
}

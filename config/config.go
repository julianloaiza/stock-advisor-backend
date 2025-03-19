package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// RecommendationFactors contiene los factores de recomendación para empresas y brokerages
type RecommendationFactors struct {
	Companies  map[string]float64 `json:"companies"`
	Brokerages map[string]float64 `json:"brokerages"`
}

// Config contiene la configuración de la aplicación.
type Config struct {
	Address               string
	DatabaseURL           string
	StockAPIURL           string
	StockAPIKey           string
	SyncMaxIterations     int
	SyncTimeout           int
	CORSAllowedOrigins    string
	RecommendationFactors *RecommendationFactors
}

// New crea una nueva instancia de Config.
func New() *Config {
	// Configuración de Viper
	setupViper()

	// Crear configuración base
	config := createBaseConfig()

	// Cargar factores de recomendación (opcional)
	loadRecommendationFactorsConfig(config)

	// Validar configuración
	if err := validateConfig(config); err != nil {
		log.Fatalf("❌ Error en la configuración: %v", err)
	}

	// Mostrar configuración (sin datos sensibles)
	logConfig(config)

	return config
}

// setupViper configura Viper y carga el archivo .env
func setupViper() {
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
	viper.SetDefault("RECOMMENDATION_FACTORS_PATH", "recommendation_factors.json")
}

// createBaseConfig crea la configuración base de la aplicación
func createBaseConfig() *Config {
	return &Config{
		Address:            viper.GetString("ADDRESS"),
		DatabaseURL:        viper.GetString("DATABASE_URL"),
		StockAPIURL:        viper.GetString("STOCK_API_URL"),
		StockAPIKey:        viper.GetString("STOCK_API_KEY"),
		SyncMaxIterations:  viper.GetInt("SYNC_MAX_ITERATIONS"),
		SyncTimeout:        viper.GetInt("SYNC_TIMEOUT"),
		CORSAllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
	}
}

// loadRecommendationFactorsConfig carga los factores de recomendación en la configuración
func loadRecommendationFactorsConfig(config *Config) {
	factorsPath := viper.GetString("RECOMMENDATION_FACTORS_PATH")
	factors, err := loadRecommendationFactors(factorsPath)
	if err != nil {
		log.Printf("ℹ️ Factores de recomendación no disponibles: %v", err)
	} else {
		config.RecommendationFactors = factors
		log.Printf("✅ Factores de recomendación cargados: %d compañías, %d brokerages",
			len(factors.Companies), len(factors.Brokerages))
	}
}

// loadRecommendationFactors carga los factores de recomendación desde un archivo JSON
func loadRecommendationFactors(path string) (*RecommendationFactors, error) {
	// Verificar si el archivo existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("archivo no encontrado")
	}

	// Leer el archivo
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo: %w", err)
	}

	// Decodificar JSON
	var factors RecommendationFactors
	if err := json.Unmarshal(data, &factors); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %w", err)
	}

	return &factors, nil
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

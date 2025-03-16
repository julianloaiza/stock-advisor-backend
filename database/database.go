package database

import (
	"log"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New crea una nueva conexión a la base de datos y migra el esquema.
func New(cfg *config.Config) *gorm.DB {
	// Configuración de GORM con nivel de log reducido
	gormConfig := &gorm.Config{
		// Usar Silent o Error para reducir drásticamente los logs
		Logger: logger.Default.LogMode(logger.Error),
	}

	// Conectar a la base de datos
	log.Println("🔌 Conectando a la base de datos...")
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}

	// Verificar la conexión con un ping simple
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Error obteniendo la conexión SQL: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Error en ping a la base de datos: %v", err)
	}

	// Auto-migrar el esquema
	log.Println("🔄 Migrando esquema de base de datos...")
	if err := db.AutoMigrate(&domain.Stock{}); err != nil {
		log.Fatalf("❌ Error en la migración: %v", err)
	}

	log.Println("✅ Conexión exitosa a la base de datos y migración completada")
	return db
}

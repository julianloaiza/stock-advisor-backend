package database

import (
	"log"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
	// "gorm.io/driver/postgres"
)

// New recibe la configuración para construir el DSN desde el archivo de configuración.
func New(cfg *config.Config) *gorm.DB {
	if cfg.SimulateDB {
		// Simulación: se salta la conexión real y se usa SQLite en memoria.
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ Error simulando la conexión a la base de datos:", err)
		}
		if err := db.AutoMigrate(&domain.Stock{}); err != nil {
			log.Fatal("❌ Error en la migración (simulada):", err)
		}
		log.Println("✅ Conexión simulada a la BD (SQLite in-memory)")
		return db
	}

	// Conexión real con CockroachDB (código comentado para uso futuro)
	/*
		db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ Error conectando a la base de datos:", err)
		}
		if err := db.AutoMigrate(&domain.Stock{}); err != nil {
			log.Fatal("❌ Error en la migración:", err)
		}
		log.Println("✅ Conexión exitosa a CockroachDB y migración completada")
		return db
	*/
	return nil
}

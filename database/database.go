package database

// New crea una nueva conexión a la base de datos y migra el esquema.
import (
	"log"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New crea una nueva conexión a la base de datos y migra el esquema.
func New(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error conectando a la base de datos:", err)
	}

	if err := db.AutoMigrate(&domain.Stock{}); err != nil {
		log.Fatal("❌ Error en la migración:", err)
	}

	log.Println("✅ Conexión exitosa a CockroachDB y migración completada")
	return db
}

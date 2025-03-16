package stocks

import (
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// ReplaceAllStocks reemplaza la data existente en la tabla Stock por la nueva data.
// Esta función es llamada desde el servicio después de haber obtenido los stocks
// de la API externa y haberlos procesado.
func (r *repository) ReplaceAllStocks(stocks []domain.Stock) error {
	log.Printf("Reemplazando todos los stocks existentes con %d nuevos registros", len(stocks))

	// Usar transacción para asegurar que ambas operaciones son atómicas
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Eliminar todos los registros existentes
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.Stock{}).Error; err != nil {
			log.Printf("Error eliminando stocks existentes: %v", err)
			return err
		}

		// Si no hay stocks para insertar, terminamos aquí
		if len(stocks) == 0 {
			log.Println("No hay stocks para insertar, tabla limpiada")
			return nil
		}

		// Insertar los nuevos registros
		// Usamos CreateInBatches para mejorar el rendimiento con grandes volúmenes de datos
		batchSize := 100
		if err := tx.CreateInBatches(&stocks, batchSize).Error; err != nil {
			log.Printf("Error insertando nuevos stocks: %v", err)
			return err
		}

		log.Printf("Sincronización completada: %d stocks reemplazados exitosamente", len(stocks))
		return nil
	})
}

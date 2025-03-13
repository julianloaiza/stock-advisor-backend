package stocks

import (
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// ReplaceAllStocks reemplaza la data existente en la tabla Stock por la nueva data.
func (r *repository) ReplaceAllStocks(stocks []domain.Stock) error {
	log.Println("Reemplazando todos los stocks existentes")
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.Stock{}).Error; err != nil {
			return err
		}
		return tx.Create(&stocks).Error
	})
}

package stocks

import (
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// Repository define las operaciones disponibles para manejar stocks.
type Repository interface {
	ReplaceAllStocks(stocks []domain.Stock) error
	GetStocks(query string, page, size int) ([]domain.Stock, int64, error)
}

type repository struct {
	db *gorm.DB
}

// New crea una nueva instancia del repositorio de stocks.
func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ReplaceAllStocks reemplaza la data existente en la tabla Stock por la nueva data,
// utilizando una transacción para asegurar atomicidad.
func (r *repository) ReplaceAllStocks(stocks []domain.Stock) error {
	log.Println("Reemplazando todos los stocks existentes")
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.Stock{}).Error; err != nil {
			return err
		}
		return tx.Create(&stocks).Error
	})
}

// GetStocks busca stocks cuyos atributos contengan el string query (case-insensitive)
// en los campos ticker, company, brokerage, action, rating_from y rating_to.
// Aplica paginación a la consulta, devolviendo además el total de registros encontrados.
func (r *repository) GetStocks(query string, page, size int) ([]domain.Stock, int64, error) {
	var stocks []domain.Stock
	var total int64
	offset := (page - 1) * size
	likeQuery := "%" + query + "%"
	dbQuery := r.db.Model(&domain.Stock{}).
		Where("ticker ILIKE ? OR company ILIKE ? OR brokerage ILIKE ? OR action ILIKE ? OR rating_from ILIKE ? OR rating_to ILIKE ?",
			likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery)

	// Contar el total de registros
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// Aplicar paginación y obtener los registros
	if err := dbQuery.Offset(offset).Limit(size).Find(&stocks).Error; err != nil {
		return nil, 0, err
	}
	return stocks, total, nil
}

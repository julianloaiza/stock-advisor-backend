package stocks

import (
	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// Repository define las operaciones disponibles para manejar stocks.
type Repository interface {
	// ReplaceAllStocks reemplaza todos los stocks en la base de datos.
	ReplaceAllStocks(stocks []domain.Stock) error

	// GetStocks obtiene los stocks filtrados, aplicando paginación en la base de datos.
	GetStocks(query string, minTargetTo, maxTargetTo float64, currency string, page, size int) ([]domain.Stock, int64, error)
}

// repository implementa la interfaz Repository.
type repository struct {
	db *gorm.DB
}

// New crea una nueva instancia del repositorio de stocks.
func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

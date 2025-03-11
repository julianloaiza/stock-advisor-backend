package stocks

import (
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// Repository define las operaciones disponibles para manejar stocks.
type Repository interface {
	SaveStock(stock domain.Stock) error
	GetStocks() ([]domain.Stock, error)
	GetStockByTicker(ticker string) (domain.Stock, error)
	UpdateStock(stock domain.Stock) error
	DeleteStock(id uint) error
}

type repository struct {
	db *gorm.DB
}

// New crea una nueva instancia del repositorio de stocks.
func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveStock(stock domain.Stock) error {
	log.Printf("Guardando stock: %s - %s", stock.Ticker, stock.Company)
	return r.db.Create(&stock).Error
}

func (r *repository) GetStocks() ([]domain.Stock, error) {
	var stocks []domain.Stock
	log.Println("Obteniendo todos los stocks")
	err := r.db.Find(&stocks).Error
	return stocks, err
}

func (r *repository) GetStockByTicker(ticker string) (domain.Stock, error) {
	var stock domain.Stock
	log.Printf("Buscando stock con ticker: %s", ticker)
	err := r.db.Where("ticker = ?", ticker).First(&stock).Error
	return stock, err
}

func (r *repository) UpdateStock(stock domain.Stock) error {
	log.Printf("Actualizando stock: %s", stock.Ticker)
	return r.db.Save(&stock).Error
}

func (r *repository) DeleteStock(id uint) error {
	log.Printf("Eliminando stock con ID: %d", id)
	return r.db.Delete(&domain.Stock{}, id).Error
}

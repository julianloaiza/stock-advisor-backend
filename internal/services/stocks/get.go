package stocks

import (
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// GetStocks maneja la búsqueda, paginación y recomendaciones.
func (s *service) GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error) {
	log.Println("Ejecutando búsqueda de stocks")

	// Obtener stocks paginados desde la base de datos
	// El repositorio ya se encarga de ordenar por RecommendScore si recommends es true
	stocks, total, err := s.repo.GetStocks(query, page, size, recommends, minTargetTo, maxTargetTo, currency)
	if err != nil {
		log.Printf("Error al obtener stocks: %v", err)
		return nil, 0, err
	}

	return stocks, total, nil
}

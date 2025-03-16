package stocks

import (
	"log"
	"sort"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// GetStocks maneja la búsqueda, paginación y recomendaciones.
func (s *service) GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error) {
	log.Println("Ejecutando búsqueda de stocks")

	// Obtener stocks paginados desde la base de datos
	allStocks, total, err := s.repo.GetStocks(query, minTargetTo, maxTargetTo, currency, page, size)
	if err != nil {
		log.Printf("Error al obtener stocks: %v", err)
		return nil, 0, err
	}

	// Si no se solicitan recomendaciones, devolvemos los resultados tal cual
	if !recommends {
		return allStocks, total, nil
	}

	// Si se solicita recomendación, reordenamos los resultados
	sortStocksByRecommendation(allStocks)

	return allStocks, total, nil
}

// sortStocksByRecommendation ordena los stocks según su puntuación de recomendación
func sortStocksByRecommendation(stocks []domain.Stock) {
	sort.Slice(stocks, func(i, j int) bool {
		return recommendationScore(stocks[i]) > recommendationScore(stocks[j])
	})
}

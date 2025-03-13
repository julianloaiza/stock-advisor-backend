package stocks

import (
	"log"
	"sort"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// GetStocks maneja la búsqueda, paginación y recomendaciones.
func (s *service) GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64) ([]domain.Stock, int64, error) {
	log.Println("Ejecutando búsqueda de stocks en el servicio")

	// Obtener stocks paginados desde la base de datos
	allStocks, total, err := s.repo.GetStocks(query, minTargetTo, maxTargetTo, page, size)
	if err != nil {
		return nil, 0, err
	}

	// Si no se solicitan recomendaciones, devolvemos los resultados tal cual
	if !recommends {
		return allStocks, total, nil
	}

	// Si se solicita recomendación, reordenamos los resultados
	sort.Slice(allStocks, func(i, j int) bool {
		return recommendationScore(allStocks[i]) > recommendationScore(allStocks[j])
	})

	return allStocks, total, nil
}

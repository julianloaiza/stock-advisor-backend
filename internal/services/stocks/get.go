package stocks

import (
	"context"
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	repo "github.com/julianloaiza/stock-advisor/internal/repositories/stocks"
)

// getStocks es la función auxiliar que delega la búsqueda de stocks en el repositorio.
func getStocks(ctx context.Context, repository repo.Repository, query string, page, size int) ([]domain.Stock, int64, error) {
	log.Println("Ejecutando búsqueda de stocks en el servicio (getStocks)")
	return repository.GetStocks(query, page, size)
}

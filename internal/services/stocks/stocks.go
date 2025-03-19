package stocks

import (
	"context"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	repo "github.com/julianloaiza/stock-advisor/internal/repositories/stocks"
	"github.com/julianloaiza/stock-advisor/internal/services/apiClient"
)

// Service define las operaciones relacionadas con stocks.
type Service interface {
	// SyncStocks sincroniza la base de datos con la API externa.
	SyncStocks(ctx context.Context, limit int) error

	// GetStocks realiza una búsqueda con query y paginación.
	GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error)
}

// service implementa la interfaz Service.
type service struct {
	repo      repo.Repository
	cfg       *config.Config
	apiClient apiClient.Client
}

// New crea una nueva instancia del servicio de stocks.
func New(repo repo.Repository, cfg *config.Config, apiClient apiClient.Client) Service {
	return &service{
		repo:      repo,
		cfg:       cfg,
		apiClient: apiClient,
	}
}

package stocks

import (
	"context"
	"time"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	repo "github.com/julianloaiza/stock-advisor/internal/repositories/stocks"
)

// Service define las operaciones relacionadas con stocks.
type Service interface {
	// SyncStocks sincroniza la base de datos con la API externa.
	SyncStocks(ctx context.Context, limit int) error
	// GetStocks realiza una búsqueda con query y paginación.
	GetStocks(ctx context.Context, query string, page, size int, recommends bool) ([]domain.Stock, int64, error)
}

type service struct {
	repo repo.Repository
	cfg  *config.Config
}

// New crea una nueva instancia del servicio de stocks.
func New(repo repo.Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

// SyncStocks utiliza la función auxiliar (definida en sync.go) para sincronizar la base de datos.
func (s *service) SyncStocks(ctx context.Context, limit int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	return syncStocks(ctx, limit, s.repo, s.cfg)
}

func (s *service) GetStocks(ctx context.Context, query string, page, size int, recommends bool) ([]domain.Stock, int64, error) {
	return getStocks(ctx, s.repo, query, page, size, recommends)
}

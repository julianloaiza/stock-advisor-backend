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
	// GetStocks obtiene los stocks desde la base de datos.
	GetStocks(ctx context.Context) ([]domain.Stock, error)
	// SyncStocks sincroniza la base de datos con la API externa.
	// El parámetro limit indica el máximo número de iteraciones (paginación) a realizar.
	SyncStocks(ctx context.Context, limit int) error
	// GetRecommendations devuelve recomendaciones de inversión.
	GetRecommendations(ctx context.Context) ([]Recommendation, error)
}

// service es la implementación concreta de Service.
type service struct {
	repo repo.Repository
	cfg  *config.Config
}

// New crea una nueva instancia del servicio de stocks.
// Se inyectan el repositorio y la configuración.
func New(repo repo.Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

// GetStocks delega la obtención de stocks al repositorio.
func (s *service) GetStocks(ctx context.Context) ([]domain.Stock, error) {
	return s.repo.GetStocks()
}

// SyncStocks utiliza la función auxiliar (definida en sync.go) para sincronizar la base de datos.
func (s *service) SyncStocks(ctx context.Context, limit int) error {
	// Obtener el timeout desde la configuración.
	timeout := s.cfg.SyncTimeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	return syncStocks(ctx, limit, s.repo, s.cfg)
}

// GetRecommendations devuelve recomendaciones (implementación simulada).
func (s *service) GetRecommendations(ctx context.Context) ([]Recommendation, error) {
	// Por el momento se devuelve una lista vacía.
	return []Recommendation{}, nil
}

package stocks

import (
	"context"
)

// Service define las operaciones relacionadas con stocks.
type Service interface {
	GetStocks(ctx context.Context) ([]Stock, error)
	SyncStocks(ctx context.Context) error
	GetRecommendations(ctx context.Context) ([]Recommendation, error)
}

// service es la implementación concreta de Service.
type service struct {
	// Aquí puedes inyectar repositorios u otras dependencias.
}

// New crea una nueva instancia del servicio de stocks.
func New() Service {
	return &service{}
}

// GetStocks delega la lógica en get.go.
func (s *service) GetStocks(ctx context.Context) ([]Stock, error) {
	return getStocks(ctx)
}

// SyncStocks delega la lógica en sync.go.
func (s *service) SyncStocks(ctx context.Context) error {
	return syncStocks(ctx)
}

// GetRecommendations delega la lógica en recommendations.go.
func (s *service) GetRecommendations(ctx context.Context) ([]Recommendation, error) {
	return getRecommendations(ctx)
}

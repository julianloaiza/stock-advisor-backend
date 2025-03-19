package apiClient

import (
	"context"
	"net/http"
	"time"

	"github.com/julianloaiza/stock-advisor/config"
)

// Client define la interfaz para comunicarse con APIs externas
type Client interface {
	// Get realiza una solicitud GET a la API
	Get(ctx context.Context, path string, params map[string]string) ([]byte, error)
}

// client implementa la interfaz Client
type client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	authHeader string
}

// New crea una nueva instancia de Cliente API basada en la configuración
func New(cfg *config.Config) Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Timeout por defecto de 30 segundos
		},
		baseURL:    cfg.StockAPIURL,
		apiKey:     cfg.StockAPIKey,
		authHeader: "X-API-KEY", // Puedes hacerlo configurable si es necesario
	}
}

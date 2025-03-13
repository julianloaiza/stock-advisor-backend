package stocks

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/julianloaiza/stock-advisor/internal/services/stocks"
)

// handler implementa la interfaz handlers.Handler.
type handler struct {
	service stocks.Service
}

// Result es el tipo que usaremos para "publicar" el handler en el grupo de handlers.
type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

// New construye el handler de stocks y lo expone como parte del grupo "handlers".
func New(service stocks.Service) Result {
	return Result{
		Handler: &handler{service: service},
	}
}

// RegisterRoutes registra las rutas de stocks.
func (h *handler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/stocks")
	group.GET("", h.GetStocks)
	group.POST("/sync", h.SyncStocks)
}

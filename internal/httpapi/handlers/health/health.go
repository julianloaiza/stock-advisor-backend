package health

import (
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// handler implementa la interfaz handlers.Handler.
type handler struct{}

// Result es el tipo para publicar el handler en el grupo de handlers.
type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

// New construye el handler de health y lo expone como parte del grupo "handlers".
func New() Result {
	return Result{
		Handler: &handler{},
	}
}

// RegisterRoutes registra las rutas de health.
func (h *handler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/health")
	group.GET("", h.HealthCheck)
}

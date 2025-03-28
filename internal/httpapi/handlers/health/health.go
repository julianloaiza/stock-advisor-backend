// internal/httpapi/handlers/health/health.go
package health

import (
	"net/http"

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
	e.GET("/health", h.HealthCheck)
}

// HealthCheck maneja el endpoint de verificación de salud.
func (h *handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

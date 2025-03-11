package profiles

import (
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/julianloaiza/stock-advisor/internal/services/profiles"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct {
	service profiles.Service
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

func New() Result {
	return Result{
		Handler: &handler{},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/profiles")
	group.GET("/:id", h.GetByID)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

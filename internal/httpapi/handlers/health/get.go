package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck maneja el endpoint de verificación de salud.
func (h *handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

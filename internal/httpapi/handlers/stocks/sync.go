package stocks

import (
	"net/http"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// SyncStocks sincroniza la base de datos con la API mediante el servicio.
func (h *handler) SyncStocks(c echo.Context) error {
	err := h.service.SyncStocks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Error:   "Error sincronizando stocks",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Data:    nil,
		Message: "Sincronización completada exitosamente",
	})
}

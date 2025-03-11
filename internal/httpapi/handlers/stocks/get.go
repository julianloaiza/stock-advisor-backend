package stocks

import (
	"net/http"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// GetStocks obtiene la lista de stocks utilizando el servicio.
func (h *handler) GetStocks(c echo.Context) error {
	stocksList, err := h.service.GetStocks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Error:   "Error obteniendo stocks",
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.APIResponse{
		Data:    stocksList,
		Message: "Consulta de acciones exitosa",
	})
}

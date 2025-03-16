package stocks

import (
	"net/http"
	"strconv"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// SyncStocks sincroniza la base de datos con la API externa.
func (h *handler) SyncStocks(c echo.Context) error {
	// Extraer y validar el límite de iteraciones (por defecto 1)
	limit := 1
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.NewError(
				http.StatusBadRequest,
				"El parámetro limit debe ser un número entero",
				err.Error(),
			))
		}
		limit = parsed
	}

	// Sincronizar los datos
	if err := h.service.SyncStocks(c.Request().Context(), limit); err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError(
			http.StatusInternalServerError,
			"Error sincronizando stocks",
			err.Error(),
		))
	}

	return c.JSON(http.StatusOK, response.NewSuccess(
		http.StatusOK,
		nil,
		"Sincronización completada exitosamente",
	))
}

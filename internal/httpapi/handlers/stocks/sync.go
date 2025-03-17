package stocks

import (
	"log"
	"net/http"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// SyncStocks sincroniza la base de datos con la API externa.
// Se espera recibir el parámetro "limit" en el body de la petición POST.
func (h *handler) SyncStocks(c echo.Context) error {
	// Usar un struct anónimo para bindear el body
	var req struct {
		Limit int `json:"limit"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewError(
			http.StatusBadRequest,
			"Error al leer el body de la petición",
			err.Error(),
		))
	}

	// Validar que se reciba un número entero positivo.
	if req.Limit <= 0 {
		return c.JSON(http.StatusBadRequest, response.NewError(
			http.StatusBadRequest,
			"El parámetro 'limit' debe ser un número entero positivo",
			"",
		))
	}
	log.Printf("Se utilizará el parámetro 'limit': %d", req.Limit)

	// Ejecutar la sincronización en el servicio.
	if err := h.service.SyncStocks(c.Request().Context(), req.Limit); err != nil {
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

package stocks

import (
	"log"
	"net/http"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// SyncRequest estructura para la solicitud de sincronización
type SyncRequest struct {
	Limit int `json:"limit" example:"5" minimum:"1"` // Número de iteraciones para la sincronización
}

// @Summary Sincronizar stocks desde fuente externa
// @Description Actualiza la base de datos con información de acciones desde un servicio externo
// @Tags stocks
// @Accept json
// @Produce json
// @Param request body SyncRequest true "Parámetros de sincronización"
// @Success 200 {object} response.APIResponse "Sincronización exitosa"
// @Failure 400 {object} response.APIResponse "Error en la solicitud"
// @Failure 500 {object} response.APIResponse "Error del servidor"
// @Router /stocks/sync [post]
func (h *handler) SyncStocks(c echo.Context) error {
	// Usar un struct para bindear el body
	var req SyncRequest
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

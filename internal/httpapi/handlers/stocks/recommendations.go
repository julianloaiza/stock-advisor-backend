package stocks

import (
	"net/http"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// GetRecommendations obtiene las recomendaciones utilizando el servicio.
func (h *handler) GetRecommendations(c echo.Context) error {
	recs, err := h.service.GetRecommendations(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Error:   "Error generando recomendaciones",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Data:    recs,
		Message: "Recomendaciones obtenidas exitosamente",
	})
}

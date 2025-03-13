package stocks

import (
	"net/http"
	"strconv"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// GetStocks obtiene la lista de stocks aplicando búsqueda y paginación.
// Parámetros de query:
//   - q: string de búsqueda (opcional)
//   - page: número de página (por defecto 1)
//   - size: tamaño de página (por defecto 10)
//   - recommends: booleano para filtrar recomendaciones (opcional)
func (h *handler) GetStocks(c echo.Context) error {
	query := c.QueryParam("query")
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("size")
	recommendsStr := c.QueryParam("recommends")

	// Valores por defecto
	page := 1
	size := 10
	recommends := false
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Code:    http.StatusBadRequest,
				Error:   "Parámetro page inválido",
				Message: "Page debe ser un entero positivo",
			})
		}
	}

	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Code:    http.StatusBadRequest,
				Error:   "Parámetro size inválido",
				Message: "Size debe ser un entero positivo",
			})
		}
	}

	if recommendsStr != "" {
		recommends, err = strconv.ParseBool(recommendsStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Code:    http.StatusBadRequest,
				Error:   "Parámetro recommends inválido",
				Message: "Recommends debe ser un booleano",
			})
		}
	}

	// Delegamos la búsqueda con paginación al servicio.
	stocksList, total, err := h.service.GetStocks(c.Request().Context(), query, page, size, recommends)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Code:    http.StatusInternalServerError,
			Error:   "Error buscando stocks",
			Message: err.Error(),
		})
	}

	paginated := response.PaginatedData{
		Content: stocksList,
		Total:   total,
		Page:    page,
		Size:    size,
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Code:    http.StatusOK,
		Data:    paginated,
		Message: "Consulta de acciones exitosa",
	})
}

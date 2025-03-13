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
//   - min_target_to: valor mínimo del target (opcional)
//   - max_target_to: valor máximo del target (opcional)
func (h *handler) GetStocks(c echo.Context) error {
	query := c.QueryParam("query")
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("size")
	recommendsStr := c.QueryParam("recommends")
	minTargetToStr := c.QueryParam("min_target_to")
	maxTargetToStr := c.QueryParam("max_target_to")

	// Valores por defecto
	page := 1
	size := 10
	recommends := false
	var minTargetTo, maxTargetTo float64
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

	if minTargetToStr != "" {
		minTargetTo, err = strconv.ParseFloat(minTargetToStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Code:    http.StatusBadRequest,
				Error:   "Parámetro min_target_to inválido",
				Message: "MinTargetTo debe ser un número",
			})
		}
	}

	if maxTargetToStr != "" {
		maxTargetTo, err = strconv.ParseFloat(maxTargetToStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Code:    http.StatusBadRequest,
				Error:   "Parámetro max_target_to inválido",
				Message: "MaxTargetTo debe ser un número",
			})
		}
	}

	// Delegamos la búsqueda con paginación al servicio.
	stocksList, total, err := h.service.GetStocks(query, page, size, recommends, minTargetTo, maxTargetTo)
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

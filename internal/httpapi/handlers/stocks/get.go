package stocks

import (
	"net/http"
	"strconv"

	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
)

// StockParams contiene los parámetros para filtrar acciones
type StockParams struct {
	Query       string
	Page        int
	Size        int
	Recommends  bool
	MinTargetTo float64
	MaxTargetTo float64
	Currency    string
}

// GetStocks
// @Summary Obtener lista de stocks
// @Description Recupera una lista filtrada y paginada de acciones bursátiles
// @Tags stocks
// @Accept json
// @Produce json
// @Param query query string false "Texto de búsqueda general (ticker, company, brokerage, etc.)"
// @Param page query int false "Número de página (por defecto: 1)" default(1)
// @Param size query int false "Registros por página (por defecto: 10)" default(10)
// @Param recommends query bool false "Ordenar por puntaje de recomendación"
// @Param minTargetTo query number false "Valor mínimo del precio objetivo"
// @Param maxTargetTo query number false "Valor máximo del precio objetivo"
// @Param currency query string false "Moneda de los precios (por defecto: USD)" default(USD)
// @Success 200 {object} response.APIResponse{data=response.PaginatedData} "Consulta de acciones exitosa"
// @Failure 400 {object} response.APIResponse "Parámetros inválidos"
// @Failure 500 {object} response.APIResponse "Error interno del servidor"
// @Router /stocks [get]
func (h *handler) GetStocks(c echo.Context) error {
	params, err := parseStockParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewError(
			http.StatusBadRequest,
			"Parámetros inválidos",
			err.Error(),
		))
	}

	// Delegamos la búsqueda con paginación al servicio
	stocksList, total, err := h.service.GetStocks(
		params.Query,
		params.Page,
		params.Size,
		params.Recommends,
		params.MinTargetTo,
		params.MaxTargetTo,
		params.Currency,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError(
			http.StatusInternalServerError,
			"Error buscando stocks",
			err.Error(),
		))
	}

	paginated := response.NewPaginated(stocksList, total, params.Page, params.Size)
	return c.JSON(http.StatusOK, response.NewSuccess(
		http.StatusOK,
		paginated,
		"Consulta de acciones exitosa",
	))
}

// parseStockParams extrae y valida los parámetros de la solicitud
func parseStockParams(c echo.Context) (StockParams, error) {
	params := StockParams{
		Query:    c.QueryParam("query"),
		Page:     1,     // Valor por defecto
		Size:     10,    // Valor por defecto
		Currency: "USD", // Valor por defecto
	}

	// Parsing de page
	if pageStr := c.QueryParam("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return params, echo.NewHTTPError(http.StatusBadRequest, "Page debe ser un entero positivo")
		}
		params.Page = page
	}

	// Parsing de size
	if sizeStr := c.QueryParam("size"); sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			return params, echo.NewHTTPError(http.StatusBadRequest, "Size debe ser un entero positivo")
		}
		params.Size = size
	}

	// Parsing de recommends
	if recommendsStr := c.QueryParam("recommends"); recommendsStr != "" {
		recommends, err := strconv.ParseBool(recommendsStr)
		if err != nil {
			return params, echo.NewHTTPError(http.StatusBadRequest, "Recommends debe ser un booleano")
		}
		params.Recommends = recommends
	}

	// Parsing de minTargetTo
	if minTargetToStr := c.QueryParam("minTargetTo"); minTargetToStr != "" {
		minTargetTo, err := strconv.ParseFloat(minTargetToStr, 64)
		if err != nil {
			return params, echo.NewHTTPError(http.StatusBadRequest, "MinTargetTo debe ser un número")
		}
		params.MinTargetTo = minTargetTo
	}

	// Parsing de maxTargetTo
	if maxTargetToStr := c.QueryParam("maxTargetTo"); maxTargetToStr != "" {
		maxTargetTo, err := strconv.ParseFloat(maxTargetToStr, 64)
		if err != nil {
			return params, echo.NewHTTPError(http.StatusBadRequest, "MaxTargetTo debe ser un número")
		}
		params.MaxTargetTo = maxTargetTo
	}

	// Parsing de currency (nuevo parámetro)
	if currencyStr := c.QueryParam("currency"); currencyStr != "" {
		params.Currency = currencyStr
	}

	return params, nil
}

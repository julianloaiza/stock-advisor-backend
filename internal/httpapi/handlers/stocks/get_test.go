package stocks

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestGetStocks_Success verifica que una solicitud válida devuelva stocks correctamente
func TestGetStocks_Success(t *testing.T) {
	// Configurar el contexto Echo con parámetros de consulta
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/stocks?query=AAPL&page=1&size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear datos de prueba
	mockStocks := []domain.Stock{
		{
			ID:         1,
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Brokerage:  "Example Broker",
			Action:     "upgraded by",
			RatingFrom: "Hold",
			RatingTo:   "Buy",
			TargetFrom: 150.0,
			TargetTo:   180.0,
			Currency:   "USD",
		},
	}

	// Crear el servicio mock
	mockService := new(mockStockService)
	mockService.On("GetStocks", "AAPL", 1, 10, false, 0.0, 0.0, "USD").Return(mockStocks, int64(1), nil)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.GetStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "Consulta de acciones exitosa", response.Message)

	// Verificar contenido paginado
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Data debería ser un mapa")

	assert.Equal(t, float64(1), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["size"])

	// Verificar que se llamó al método del servicio con los parámetros correctos
	mockService.AssertExpectations(t)
}

// TestGetStocks_WithRecommendation verifica que el parámetro recommends funcione
func TestGetStocks_WithRecommendation(t *testing.T) {
	// Configurar el contexto Echo con parámetros de consulta incluyendo recommends=true
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/stocks?query=TECH&page=1&size=10&recommends=true", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear datos de prueba
	mockStocks := []domain.Stock{
		{
			ID:         1,
			Ticker:     "TSLA",
			Company:    "Tesla Inc.",
			Brokerage:  "Example Broker",
			Action:     "target raised by",
			RatingFrom: "Buy",
			RatingTo:   "Buy",
			TargetFrom: 800.0,
			TargetTo:   900.0,
			Currency:   "USD",
		},
	}

	// Crear el servicio mock
	mockService := new(mockStockService)
	mockService.On("GetStocks", "TECH", 1, 10, true, 0.0, 0.0, "USD").Return(mockStocks, int64(1), nil)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.GetStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.Code)

	// Verificar que se llamó al método del servicio con recommends=true
	mockService.AssertExpectations(t)
}

// TestGetStocks_WithPriceFilter verifica que los filtros de precio funcionen
func TestGetStocks_WithPriceFilter(t *testing.T) {
	// Configurar el contexto Echo con parámetros de filtro de precio
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/stocks?minTargetTo=100&maxTargetTo=200&currency=EUR", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear datos de prueba
	mockStocks := []domain.Stock{
		{
			ID:         1,
			Ticker:     "BMW",
			Company:    "Bayerische Motoren Werke AG",
			Brokerage:  "European Broker",
			Action:     "reiterated by",
			RatingFrom: "Hold",
			RatingTo:   "Hold",
			TargetFrom: 120.0,
			TargetTo:   130.0,
			Currency:   "EUR",
		},
	}

	// Crear el servicio mock
	mockService := new(mockStockService)
	mockService.On("GetStocks", "", 1, 10, false, 100.0, 200.0, "EUR").Return(mockStocks, int64(1), nil)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.GetStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verificar que se llamó al método del servicio con los filtros correctos
	mockService.AssertExpectations(t)
}

// TestGetStocks_InvalidPage verifica que se maneje adecuadamente una página inválida
func TestGetStocks_InvalidPage(t *testing.T) {
	// Configurar el contexto Echo con parámetro de página inválido
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/stocks?page=invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock (no debería ser llamado)
	mockService := new(mockStockService)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.GetStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "Parámetros inválidos")

	// Verificar que NO se llamó al método del servicio
	mockService.AssertNotCalled(t, "GetStocks")
}

// TestGetStocks_ServiceError verifica que se maneje adecuadamente un error del servicio
func TestGetStocks_ServiceError(t *testing.T) {
	// Configurar el contexto Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock que devolverá un error
	mockService := new(mockStockService)
	expectedError := errors.New("error de base de datos")
	mockService.On("GetStocks", "", 1, 10, false, 0.0, 0.0, "USD").Return([]domain.Stock{}, int64(0), expectedError)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.GetStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Error buscando stocks", response.Message)
	assert.Equal(t, expectedError.Error(), response.Error)

	// Verificar que se llamó al método del servicio
	mockService.AssertExpectations(t)
}

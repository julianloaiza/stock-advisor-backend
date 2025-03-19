package stocks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Servicio mock para pruebas
type mockStockService struct {
	mock.Mock
}

func (m *mockStockService) GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error) {
	args := m.Called(query, page, size, recommends, minTargetTo, maxTargetTo, currency)
	return args.Get(0).([]domain.Stock), args.Get(1).(int64), args.Error(2)
}

func (m *mockStockService) SyncStocks(ctx context.Context, limit int) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

// TestSyncStocks_Success verifica que la sincronización exitosa devuelva un código 200
func TestSyncStocks_Success(t *testing.T) {
	// Configurar el contexto Echo y la solicitud
	e := echo.New()
	requestBody := map[string]int{"limit": 5}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/stocks/sync", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock y configurar comportamiento esperado
	mockService := new(mockStockService)
	mockService.On("SyncStocks", mock.Anything, 5).Return(nil)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.SyncStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "Sincronización completada exitosamente", response.Message)
	assert.Nil(t, response.Data)
	assert.Empty(t, response.Error)

	// Verificar que se llamó al método del servicio con los parámetros correctos
	mockService.AssertExpectations(t)
}

// TestSyncStocks_InvalidLimit verifica que un límite inválido devuelva un error 400
func TestSyncStocks_InvalidLimit(t *testing.T) {
	// Configurar el contexto Echo y la solicitud con límite inválido
	e := echo.New()
	requestBody := map[string]int{"limit": 0} // Límite inválido
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/stocks/sync", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock (no debería ser llamado)
	mockService := new(mockStockService)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.SyncStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "El parámetro 'limit' debe ser un número entero positivo")
	assert.Nil(t, response.Data)

	// Verificar que NO se llamó al método del servicio
	mockService.AssertNotCalled(t, "SyncStocks")
}

// TestSyncStocks_InvalidBody verifica que un cuerpo de solicitud inválido devuelva un error 400
func TestSyncStocks_InvalidBody(t *testing.T) {
	// Configurar el contexto Echo y la solicitud con JSON inválido
	e := echo.New()
	invalidJSON := []byte(`{"limit": "not-a-number"}`)
	req := httptest.NewRequest(http.MethodPost, "/stocks/sync", bytes.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock (no debería ser llamado)
	mockService := new(mockStockService)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.SyncStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "Error al leer el body de la petición")
	assert.Nil(t, response.Data)

	// Verificar que NO se llamó al método del servicio
	mockService.AssertNotCalled(t, "SyncStocks")
}

// TestSyncStocks_ServiceError verifica que un error del servicio devuelva un error 500
func TestSyncStocks_ServiceError(t *testing.T) {
	// Configurar el contexto Echo y la solicitud
	e := echo.New()
	requestBody := map[string]int{"limit": 5}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/stocks/sync", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock y configurar un error
	mockService := new(mockStockService)
	expectedError := errors.New("error de sincronización")
	mockService.On("SyncStocks", mock.Anything, 5).Return(expectedError)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.SyncStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Verificar respuesta JSON
	var response response.APIResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Error sincronizando stocks", response.Message)
	assert.Equal(t, expectedError.Error(), response.Error)
	assert.Nil(t, response.Data)

	// Verificar que se llamó al método del servicio con los parámetros correctos
	mockService.AssertExpectations(t)
}

// TestSyncStocks_EmptyBody verifica que una solicitud sin cuerpo devuelva un error 400
func TestSyncStocks_EmptyBody(t *testing.T) {
	// Configurar el contexto Echo y la solicitud sin cuerpo
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/stocks/sync", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Crear el servicio mock (no debería ser llamado)
	mockService := new(mockStockService)

	// Crear el handler con el servicio mock
	h := &handler{service: mockService}

	// Ejecutar el handler
	err := h.SyncStocks(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verificar que NO se llamó al método del servicio
	mockService.AssertNotCalled(t, "SyncStocks")
}

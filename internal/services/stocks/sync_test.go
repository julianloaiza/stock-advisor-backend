package stocks

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio de stocks
type MockRepository struct {
	mock.Mock
}

// ReplaceAllStocks implementa el método del repositorio
func (m *MockRepository) ReplaceAllStocks(stocks []domain.Stock) error {
	args := m.Called(stocks)
	return args.Error(0)
}

// GetStocks implementa el método del repositorio para cumplir con la interfaz
func (m *MockRepository) GetStocks(query string, minTargetTo, maxTargetTo float64, currency string, page, size int) ([]domain.Stock, int64, error) {
	args := m.Called(query, minTargetTo, maxTargetTo, currency, page, size)
	return args.Get(0).([]domain.Stock), args.Get(1).(int64), args.Error(2)
}

// createMockConfig crea una configuración mock para las pruebas
func createMockConfig(apiURL, apiKey string, maxIterations, timeout int) *config.Config {
	return &config.Config{
		StockAPIURL:       apiURL,
		StockAPIKey:       apiKey,
		SyncMaxIterations: maxIterations,
		SyncTimeout:       timeout,
		// Otros campos requeridos pueden dejarse con valores por defecto
		Address:            ":8080",
		DatabaseURL:        "mock-db-url",
		CORSAllowedOrigins: "*",
	}
}

// TestSyncStocks_Success prueba que la sincronización se realice correctamente
func TestSyncStocks_Success(t *testing.T) {
	// Crear servidor HTTP mock para simular API externa
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar que la solicitud tenga Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Contains(t, authHeader, "Bearer ", "La solicitud debe incluir un token Bearer")

		// Simular respuesta JSON de la API externa
		jsonResponse := `{
			"items": [
				{
					"ticker": "AAPL",
					"company": "Apple Inc.",
					"brokerage": "Example Brokerage",
					"action": "target raised by",
					"rating_from": "Buy",
					"rating_to": "Strong-Buy",
					"target_from": "150.00",
					"target_to": "180.00",
					"currency": "USD"
				}
			],
			"next_page": ""
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	// Crear mock del repositorio
	mockRepo := new(MockRepository)

	// Configurar expectativa: el repositorio recibirá una llamada a ReplaceAllStocks
	mockRepo.On("ReplaceAllStocks", mock.Anything).Return(nil)

	// Crear configuración mock
	mockCfg := createMockConfig(server.URL, "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock y la configuración mock
	service := NewService(mockRepo, mockCfg)

	// Ejecutar el método a probar
	err := service.SyncStocks(context.Background(), 1)

	// Verificar que no hay error
	assert.NoError(t, err, "La sincronización debería ser exitosa")

	// Verificar que el repositorio fue llamado con un slice de stocks
	mockRepo.AssertExpectations(t)

	// Verificar que se llamó a ReplaceAllStocks con al menos un stock
	calls := mockRepo.Calls
	assert.GreaterOrEqual(t, len(calls), 1, "Debería haber al menos una llamada al repositorio")

	// Verificar el primer argumento de la llamada a ReplaceAllStocks
	if len(calls) > 0 {
		stocks, ok := calls[0].Arguments[0].([]domain.Stock)
		assert.True(t, ok, "El primer argumento debería ser un slice de stocks")
		assert.GreaterOrEqual(t, len(stocks), 1, "Debería haber al menos un stock para reemplazar")

		// Verificar datos del stock si hay al menos uno
		if len(stocks) > 0 {
			assert.Equal(t, "AAPL", stocks[0].Ticker, "El ticker debería ser AAPL")
			assert.Equal(t, "Apple Inc.", stocks[0].Company, "La compañía debería ser Apple Inc.")
			assert.Equal(t, 180.0, stocks[0].TargetTo, "El target_to debería ser 180.0")
		}
	}
}

// TestSyncStocks_ExternalAPIError prueba que se maneje correctamente un error de la API externa
func TestSyncStocks_ExternalAPIError(t *testing.T) {
	// Crear servidor HTTP mock que responde con error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
	}))
	defer server.Close()

	// Crear mock del repositorio
	mockRepo := new(MockRepository)
	// No esperamos que se llame a ReplaceAllStocks porque la API externa fallará

	// Crear configuración mock
	mockCfg := createMockConfig(server.URL, "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock y la configuración mock
	service := NewService(mockRepo, mockCfg)

	// Ejecutar el método a probar
	err := service.SyncStocks(context.Background(), 1)

	// Verificar que hay un error
	assert.Error(t, err, "La sincronización debería fallar")
	assert.Contains(t, err.Error(), "status code inesperado", "El error debería mencionar el status code")
}

// TestSyncStocks_RepositoryError prueba que se maneje correctamente un error del repositorio
func TestSyncStocks_RepositoryError(t *testing.T) {
	// Crear servidor HTTP mock para simular API externa
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse := `{
			"items": [
				{
					"ticker": "AAPL",
					"company": "Apple Inc.",
					"brokerage": "Example Brokerage",
					"action": "target raised by",
					"rating_from": "Buy",
					"rating_to": "Strong-Buy",
					"target_from": "150.00",
					"target_to": "180.00",
					"currency": "USD"
				}
			],
			"next_page": ""
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	// Crear mock del repositorio que devuelve error
	mockRepo := new(MockRepository)

	// Configurar expectativa: el repositorio devolverá un error
	expectedError := errors.New("error al reemplazar stocks")
	mockRepo.On("ReplaceAllStocks", mock.Anything).Return(expectedError)

	// Crear configuración mock
	mockCfg := createMockConfig(server.URL, "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock y la configuración mock
	service := NewService(mockRepo, mockCfg)

	// Ejecutar el método a probar
	err := service.SyncStocks(context.Background(), 1)

	// Verificar que hay un error y es el esperado
	assert.Error(t, err, "La sincronización debería fallar")
	assert.Contains(t, err.Error(), "error reemplazando stocks", "El error debería mencionar el problema en el repositorio")

	// Verificar que el repositorio fue llamado
	mockRepo.AssertExpectations(t)
}

// TestSyncStocks_InvalidLimit prueba que se maneje correctamente un límite inválido
func TestSyncStocks_InvalidLimit(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(MockRepository)

	// No configuramos expectativas porque no debería llamarse al repositorio

	// Crear configuración mock
	mockCfg := createMockConfig("http://example.com", "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock y la configuración mock
	service := NewService(mockRepo, mockCfg)

	// Probar con límite negativo
	ctx := context.Background()

	// Ejecutar método con límite inválido y verificar que use el valor por defecto (1)
	err := service.SyncStocks(ctx, -5)

	// Debería fallar porque no hay un servidor HTTP real para responder
	// Pero lo importante es verificar que la función no falla por el límite inválido
	assert.Error(t, err, "Debería fallar por falta de servidor HTTP, no por el límite inválido")

	// El error debería estar relacionado con la conexión HTTP, no con el límite
	assert.NotContains(t, err.Error(), "límite inválido", "El error no debería estar relacionado con el límite")
}

// NewService crea un nuevo servicio para las pruebas
func NewService(repo Repository, cfg *config.Config) *service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

// Repository define la interfaz del repositorio
type Repository interface {
	GetStocks(query string, minTargetTo, maxTargetTo float64, currency string, page, size int) ([]domain.Stock, int64, error)
	ReplaceAllStocks(stocks []domain.Stock) error
}

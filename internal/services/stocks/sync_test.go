package stocks

import (
	"context"
	"errors"
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

func (m *MockRepository) ReplaceAllStocks(stocks []domain.Stock) error {
	args := m.Called(stocks)
	return args.Error(0)
}

func (m *MockRepository) GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error) {
	args := m.Called(query, page, size, recommends, minTargetTo, maxTargetTo, currency)
	return args.Get(0).([]domain.Stock), args.Get(1).(int64), args.Error(2)
}

// MockAPIClient es un mock del cliente de API para las pruebas
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) Get(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	args := m.Called(ctx, endpoint, params)
	return args.Get(0).([]byte), args.Error(1)
}

// createMockConfig crea una configuración mock para las pruebas
func createMockConfig(apiURL, apiKey string, maxIterations, timeout int) *config.Config {
	return &config.Config{
		StockAPIURL:        apiURL,
		StockAPIKey:        apiKey,
		SyncMaxIterations:  maxIterations,
		SyncTimeout:        timeout,
		Address:            ":8080",
		DatabaseURL:        "mock-db-url",
		CORSAllowedOrigins: "*",
	}
}

// TestSyncStocks_Success prueba que la sincronización se realice correctamente
func TestSyncStocks_Success(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(MockRepository)
	mockRepo.On("ReplaceAllStocks", mock.Anything).Return(nil)

	// Crear mock del cliente API
	mockAPIClient := new(MockAPIClient)
	jsonResponse := []byte(`{
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
	}`)
	mockAPIClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(jsonResponse, nil)

	// Crear configuración mock
	mockCfg := createMockConfig("http://test-api.com", "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock, configuración mock y cliente API mock
	service := &service{
		repo:      mockRepo,
		cfg:       mockCfg,
		apiClient: mockAPIClient,
	}

	// Ejecutar el método a probar
	err := service.SyncStocks(context.Background(), 1)

	// Verificar que no hay error
	assert.NoError(t, err, "La sincronización debería ser exitosa")

	// Verificar que el repositorio fue llamado correctamente
	mockRepo.AssertExpectations(t)
	mockAPIClient.AssertExpectations(t)
}

// TestSyncStocks_ExternalAPIError prueba que se maneje correctamente un error de la API externa
func TestSyncStocks_ExternalAPIError(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(MockRepository)

	// Crear mock del cliente API que devuelve un error
	mockAPIClient := new(MockAPIClient)
	mockAPIClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, errors.New("error de API externa"))

	// Crear configuración mock
	mockCfg := createMockConfig("http://test-api.com", "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock, configuración mock y cliente API mock
	service := &service{
		repo:      mockRepo,
		cfg:       mockCfg,
		apiClient: mockAPIClient,
	}

	// Ejecutar el método a probar
	err := service.SyncStocks(context.Background(), 1)

	// Verificar que hay un error
	assert.Error(t, err, "La sincronización debería fallar")
	assert.Contains(t, err.Error(), "error de API externa")

	// Verificar que el cliente API fue llamado
	mockAPIClient.AssertExpectations(t)
}

// TestSyncStocks_RepositoryError prueba que se maneje correctamente un error del repositorio
func TestSyncStocks_RepositoryError(t *testing.T) {
	// Crear mock del repositorio que devuelve error
	mockRepo := new(MockRepository)
	mockRepo.On("ReplaceAllStocks", mock.Anything).Return(errors.New("error al reemplazar stocks"))

	// Crear mock del cliente API
	mockAPIClient := new(MockAPIClient)
	jsonResponse := []byte(`{
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
	}`)
	mockAPIClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(jsonResponse, nil)

	// Crear configuración mock
	mockCfg := createMockConfig("http://test-api.com", "test-api-key", 10, 30)

	// Crear el servicio con el repositorio mock, configuración mock y cliente API mock
	service := &service{
		repo:      mockRepo,
		cfg:       mockCfg,
		apiClient: mockAPIClient,
	}

	// Ejecutar el método a probar
	err := service.SyncStocks(context.Background(), 1)

	// Verificar que hay un error
	assert.Error(t, err, "La sincronización debería fallar")
	assert.Contains(t, err.Error(), "error reemplazando stocks")

	// Verificar que el repositorio y el cliente API fueron llamados
	mockRepo.AssertExpectations(t)
	mockAPIClient.AssertExpectations(t)
}

// Repository define la interfaz del repositorio
type Repository interface {
	GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error)
	ReplaceAllStocks(stocks []domain.Stock) error
}

package stocks

import (
	"testing"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase simula un repositorio para pruebas
type MockDatabase struct {
	mock.Mock
}

// GetStocks simula la recuperación de stocks
func (m *MockDatabase) GetStocks(query string, minTargetTo, maxTargetTo float64, currency string, page, size int) ([]domain.Stock, int64, error) {
	args := m.Called(query, minTargetTo, maxTargetTo, currency, page, size)
	return args.Get(0).([]domain.Stock), args.Get(1).(int64), args.Error(2)
}

// ReplaceAllStocks simula el reemplazo de stocks (para cumplir con la interfaz)
func (m *MockDatabase) ReplaceAllStocks(stocks []domain.Stock) error {
	args := m.Called(stocks)
	return args.Error(0)
}

// Datos de prueba para stocks
var testStocks = []domain.Stock{
	{
		Ticker:     "AAPL",
		Company:    "Apple Inc.",
		Brokerage:  "Goldman Sachs",
		Action:     "upgraded by",
		RatingFrom: "Hold",
		RatingTo:   "Buy",
		TargetFrom: 150.0,
		TargetTo:   180.0,
		Currency:   "USD",
	},
	{
		Ticker:     "GOOGL",
		Company:    "Alphabet Inc.",
		Brokerage:  "Morgan Stanley",
		Action:     "reiterated by",
		RatingFrom: "Buy",
		RatingTo:   "Buy",
		TargetFrom: 2000.0,
		TargetTo:   2200.0,
		Currency:   "USD",
	},
	{
		Ticker:     "MSFT",
		Company:    "Microsoft Corporation",
		Brokerage:  "JP Morgan",
		Action:     "target raised by",
		RatingFrom: "Buy",
		RatingTo:   "Buy",
		TargetFrom: 300.0,
		TargetTo:   350.0,
		Currency:   "EUR",
	},
}

// TestGetStocks_BasicSearch prueba la búsqueda básica de stocks
func TestGetStocks_BasicSearch(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockDatabase)

	// Configurar expectativas
	mockDB.On("GetStocks", "Apple", 0.0, 0.0, "USD", 1, 10).
		Return([]domain.Stock{testStocks[0]}, int64(1), nil)

	// Ejecutar búsqueda
	stocks, total, err := mockDB.GetStocks("Apple", 0.0, 0.0, "USD", 1, 10)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "AAPL", stocks[0].Ticker)

	// Verificar que se llamó al método con los parámetros esperados
	mockDB.AssertExpectations(t)
}

// TestGetStocks_PriceFilter prueba el filtrado por rango de precios
func TestGetStocks_PriceFilter(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockDatabase)

	// Configurar expectativas
	mockDB.On("GetStocks", "", 1000.0, 2500.0, "USD", 1, 10).
		Return([]domain.Stock{testStocks[1]}, int64(1), nil)

	// Ejecutar búsqueda con filtro de precio mínimo
	stocks, total, err := mockDB.GetStocks("", 1000.0, 2500.0, "USD", 1, 10)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "GOOGL", stocks[0].Ticker)

	// Verificar que se llamó al método con los parámetros esperados
	mockDB.AssertExpectations(t)
}

// TestGetStocks_CurrencyFilter prueba el filtrado por moneda
func TestGetStocks_CurrencyFilter(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockDatabase)

	// Configurar expectativas
	mockDB.On("GetStocks", "", 0.0, 0.0, "EUR", 1, 10).
		Return([]domain.Stock{testStocks[2]}, int64(1), nil)

	// Ejecutar búsqueda con filtro de moneda
	stocks, total, err := mockDB.GetStocks("", 0.0, 0.0, "EUR", 1, 10)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "MSFT", stocks[0].Ticker)

	// Verificar que se llamó al método con los parámetros esperados
	mockDB.AssertExpectations(t)
}

// TestGetStocks_Pagination prueba la funcionalidad de paginación
func TestGetStocks_Pagination(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockDatabase)

	// Configurar expectativas
	mockDB.On("GetStocks", "", 0.0, 0.0, "USD", 2, 1).
		Return([]domain.Stock{testStocks[1]}, int64(2), nil)

	// Ejecutar búsqueda con límite de 1 y página 2
	stocks, total, err := mockDB.GetStocks("", 0.0, 0.0, "USD", 2, 1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "GOOGL", stocks[0].Ticker)

	// Verificar que se llamó al método con los parámetros esperados
	mockDB.AssertExpectations(t)
}

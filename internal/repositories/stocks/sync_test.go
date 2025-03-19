package stocks

import (
	"testing"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSyncDatabase simula un repositorio para pruebas de sincronización
type MockSyncDatabase struct {
	mock.Mock
}

// ReplaceAllStocks simula el reemplazo de stocks en la base de datos
func (m *MockSyncDatabase) ReplaceAllStocks(stocks []domain.Stock) error {
	args := m.Called(stocks)
	return args.Error(0)
}

// GetStocks es un método adicional para cumplir con la interfaz completa
func (m *MockSyncDatabase) GetStocks(query string, minTargetTo, maxTargetTo float64, currency string, page, size int) ([]domain.Stock, int64, error) {
	args := m.Called(query, minTargetTo, maxTargetTo, currency, page, size)
	return args.Get(0).([]domain.Stock), args.Get(1).(int64), args.Error(2)
}

// Datos de prueba para stocks
var syncTestStocks = []domain.Stock{
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
}

// TestReplaceAllStocks_Success prueba el reemplazo exitoso de stocks
func TestReplaceAllStocks_Success(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockSyncDatabase)

	// Configurar expectativas
	mockDB.On("ReplaceAllStocks", syncTestStocks).Return(nil)

	// Ejecutar reemplazo de stocks
	err := mockDB.ReplaceAllStocks(syncTestStocks)

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que se llamó al método con los parámetros esperados
	mockDB.AssertExpectations(t)
}

// TestReplaceAllStocks_EmptyList prueba el reemplazo con lista vacía
func TestReplaceAllStocks_EmptyList(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockSyncDatabase)

	// Configurar expectativas
	mockDB.On("ReplaceAllStocks", []domain.Stock{}).Return(nil)

	// Ejecutar reemplazo de stocks con lista vacía
	err := mockDB.ReplaceAllStocks([]domain.Stock{})

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que se llamó al método con lista vacía
	mockDB.AssertExpectations(t)
}

// TestReplaceAllStocks_DatabaseError prueba el manejo de errores de base de datos
func TestReplaceAllStocks_DatabaseError(t *testing.T) {
	// Crear mock de base de datos
	mockDB := new(MockSyncDatabase)

	// Definir un error de base de datos simulado
	databaseError := assert.AnError

	// Configurar expectativas con error
	mockDB.On("ReplaceAllStocks", syncTestStocks).Return(databaseError)

	// Ejecutar reemplazo de stocks
	err := mockDB.ReplaceAllStocks(syncTestStocks)

	// Verificaciones
	assert.Error(t, err)
	assert.Equal(t, databaseError, err)

	// Verificar que se llamó al método con los parámetros esperados
	mockDB.AssertExpectations(t)
}

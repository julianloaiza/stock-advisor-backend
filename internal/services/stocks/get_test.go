package stocks

import (
	"testing"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Repositorio mock para pruebas
type mockStockRepository struct {
	mock.Mock
}

func (m *mockStockRepository) GetStocks(query string, minTargetTo, maxTargetTo float64, currency string, page, size int) ([]domain.Stock, int64, error) {
	args := m.Called(query, minTargetTo, maxTargetTo, currency, page, size)
	return args.Get(0).([]domain.Stock), args.Get(1).(int64), args.Error(2)
}

func (m *mockStockRepository) ReplaceAllStocks(stocks []domain.Stock) error {
	args := m.Called(stocks)
	return args.Error(0)
}

// TestGetStocks_BasicQuery prueba una consulta básica sin recomendaciones
func TestGetStocks_BasicQuery(t *testing.T) {
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
		{
			ID:         2,
			Ticker:     "MSFT",
			Company:    "Microsoft Corporation",
			Brokerage:  "Another Broker",
			Action:     "reiterated by",
			RatingFrom: "Buy",
			RatingTo:   "Buy",
			TargetFrom: 300.0,
			TargetTo:   350.0,
			Currency:   "USD",
		},
	}

	// Crear repositorio mock
	mockRepo := new(mockStockRepository)
	mockRepo.On("GetStocks", "tech", 0.0, 0.0, "USD", 1, 10).Return(mockStocks, int64(2), nil)

	// Crear el servicio con el repositorio mock
	s := &service{repo: mockRepo}

	// Ejecutar función del servicio
	result, total, err := s.GetStocks("tech", 1, 10, false, 0.0, 0.0, "USD")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "AAPL", result[0].Ticker)
	assert.Equal(t, "MSFT", result[1].Ticker)

	// Verificar que se llamó al método del repositorio con los parámetros correctos
	mockRepo.AssertExpectations(t)
}

// TestGetStocks_WithRecommendations prueba la funcionalidad de recomendaciones
func TestGetStocks_WithRecommendations(t *testing.T) {
	// Crear datos de prueba que incluyen stocks con diferentes puntuaciones potenciales
	mockStocks := []domain.Stock{
		{
			// Stock con puntuación baja (debe aparecer después en las recomendaciones)
			ID:         1,
			Ticker:     "LOW",
			Company:    "Low Score Inc.",
			Brokerage:  "Broker A",
			Action:     "reiterated by",
			RatingFrom: "Hold",
			RatingTo:   "Hold",
			TargetFrom: 100.0,
			TargetTo:   101.0, // Pequeña diferencia
			Currency:   "USD",
		},
		{
			// Stock con puntuación alta (debe aparecer primero en las recomendaciones)
			ID:         2,
			Ticker:     "HIGH",
			Company:    "High Score Corp.",
			Brokerage:  "Broker B",
			Action:     "upgraded by",
			RatingFrom: "Hold",
			RatingTo:   "Strong-Buy",
			TargetFrom: 200.0,
			TargetTo:   300.0, // Gran diferencia
			Currency:   "USD",
		},
	}

	// Crear repositorio mock
	mockRepo := new(mockStockRepository)
	mockRepo.On("GetStocks", "invest", 0.0, 0.0, "USD", 1, 10).Return(mockStocks, int64(2), nil)

	// Crear el servicio con el repositorio mock
	s := &service{repo: mockRepo}

	// Ejecutar función del servicio con recommends = true
	result, total, err := s.GetStocks("invest", 1, 10, true, 0.0, 0.0, "USD")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))

	// Verificar que el resultado está ordenado por recomendación
	// El stock con puntuación más alta debe aparecer primero
	assert.Equal(t, "HIGH", result[0].Ticker)
	assert.Equal(t, "LOW", result[1].Ticker)

	// Verificar que se llamó al método del repositorio con los parámetros correctos
	mockRepo.AssertExpectations(t)
}

// TestGetStocks_WithFilters prueba el uso de filtros de precio y moneda
func TestGetStocks_WithFilters(t *testing.T) {
	// Crear datos de prueba
	mockStocks := []domain.Stock{
		{
			ID:         1,
			Ticker:     "EUR",
			Company:    "Euro Stock",
			Brokerage:  "European Broker",
			Action:     "upgraded by",
			RatingFrom: "Hold",
			RatingTo:   "Buy",
			TargetFrom: 50.0,
			TargetTo:   75.0,
			Currency:   "EUR",
		},
	}

	// Crear repositorio mock
	mockRepo := new(mockStockRepository)
	mockRepo.On("GetStocks", "", 50.0, 100.0, "EUR", 1, 20).Return(mockStocks, int64(1), nil)

	// Crear el servicio con el repositorio mock
	s := &service{repo: mockRepo}

	// Ejecutar función del servicio con filtros
	result, total, err := s.GetStocks("", 1, 20, false, 50.0, 100.0, "EUR")

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "EUR", result[0].Ticker)
	assert.Equal(t, "EUR", result[0].Currency)

	// Verificar que se llamó al método del repositorio con los parámetros correctos
	mockRepo.AssertExpectations(t)
}

// TestRecommendationScore prueba que el algoritmo de recomendación funciona correctamente
func TestRecommendationScore(t *testing.T) {
	testCases := []struct {
		name     string
		stock    domain.Stock
		expected bool // true si debe tener una puntuación más alta que el caso siguiente
	}{
		{
			name: "Alta puntuación - Gran aumento de precio y mejora de calificación",
			stock: domain.Stock{
				Ticker:     "HIGH",
				Action:     "upgraded by",
				RatingFrom: "Hold",
				RatingTo:   "Strong-Buy",
				TargetFrom: 100.0,
				TargetTo:   200.0, // 100% de aumento
			},
			expected: true,
		},
		{
			name: "Puntuación media - Aumento moderado",
			stock: domain.Stock{
				Ticker:     "MED",
				Action:     "target raised by",
				RatingFrom: "Buy",
				RatingTo:   "Buy",
				TargetFrom: 100.0,
				TargetTo:   130.0, // 30% de aumento
			},
			expected: true,
		},
		{
			name: "Baja puntuación - Pequeño aumento",
			stock: domain.Stock{
				Ticker:     "LOW",
				Action:     "reiterated by",
				RatingFrom: "Hold",
				RatingTo:   "Hold",
				TargetFrom: 100.0,
				TargetTo:   105.0, // 5% de aumento
			},
			expected: false,
		},
	}

	// Comparar puntuaciones entre casos consecutivos
	for i := 0; i < len(testCases)-1; i++ {
		current := recommendationScore(testCases[i].stock)
		next := recommendationScore(testCases[i+1].stock)

		if testCases[i].expected {
			assert.Greater(t, current, next,
				"La puntuación de %s debería ser mayor que %s",
				testCases[i].name, testCases[i+1].name)
		} else {
			assert.LessOrEqual(t, current, next,
				"La puntuación de %s debería ser menor o igual que %s",
				testCases[i].name, testCases[i+1].name)
		}
	}
}

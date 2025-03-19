package stocks

import (
	"testing"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestCalculateRecommendationScore verifica el algoritmo de puntuación de recomendaciones
func TestCalculateRecommendationScore(t *testing.T) {
	// Casos de prueba para diferentes escenarios de puntuación
	testCases := []struct {
		name     string
		stock    domain.Stock
		expected float64
	}{
		{
			name: "Alta Puntuación - Aumento Significativo de Precio y Mejora de Calificación",
			stock: domain.Stock{
				Ticker:     "BEST",
				Action:     "upgraded by",
				RatingFrom: "Hold",
				RatingTo:   "Strong-Buy",
				TargetFrom: 100.0,
				TargetTo:   200.0, // Aumento del 100%
			},
			expected: 132, // Puntuación esperada alta
		},
		{
			name: "Puntuación Media - Aumento Moderado",
			stock: domain.Stock{
				Ticker:     "GOOD",
				Action:     "target raised by",
				RatingFrom: "Buy",
				RatingTo:   "Buy",
				TargetFrom: 100.0,
				TargetTo:   130.0, // Aumento del 30%
			},
			expected: 50, // Puntuación esperada media
		},
		{
			name: "Puntuación Baja - Aumento Mínimo",
			stock: domain.Stock{
				Ticker:     "LOW",
				Action:     "reiterated by",
				RatingFrom: "Hold",
				RatingTo:   "Hold",
				TargetFrom: 100.0,
				TargetTo:   105.0, // Aumento del 5%
			},
			expected: 7, // Puntuación esperada baja
		},
	}

	// Iterar sobre los casos de prueba
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calcular la puntuación de recomendación
			score := recommendationScore(tc.stock)

			// Verificar que la puntuación sea exactamente la esperada
			assert.Equal(t, tc.expected, score,
				"La puntuación calculada debe coincidir con la esperada para %s", tc.name)
		})
	}
}

// TestRecommendationScoreSorting verifica que las acciones se ordenen correctamente por puntuación de recomendación
func TestRecommendationScoreSorting(t *testing.T) {
	// Crear un slice de acciones con diferentes puntuaciones
	stocks := []domain.Stock{
		{
			Ticker:     "LOW",
			Action:     "reiterated by",
			RatingFrom: "Hold",
			RatingTo:   "Hold",
			TargetFrom: 100.0,
			TargetTo:   105.0,
		},
		{
			Ticker:     "HIGH",
			Action:     "upgraded by",
			RatingFrom: "Hold",
			RatingTo:   "Strong-Buy",
			TargetFrom: 100.0,
			TargetTo:   200.0,
		},
	}

	// Ordenar las acciones por puntuación de recomendación
	sortStocksByRecommendation(stocks)

	// Verificar que las acciones estén ordenadas en orden descendente de puntuación
	assert.Equal(t, "HIGH", stocks[0].Ticker, "La acción con mayor puntuación debe estar primero")
	assert.Equal(t, "LOW", stocks[1].Ticker, "La acción con menor puntuación debe estar al final")
}

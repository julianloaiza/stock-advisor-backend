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
			expected: 36.125, // Puntuación esperada alta (valor real)
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
			expected: 21.15, // Puntuación esperada media (valor real)
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
			expected: 3.8, // Puntuación esperada baja (valor real)
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

// TestRecommendationScoreComparison verifica que las acciones con mejores características tengan mayores puntuaciones
func TestRecommendationScoreComparison(t *testing.T) {
	testCases := []struct {
		name  string
		stock domain.Stock
	}{
		{
			name: "Alta Puntuación",
			stock: domain.Stock{
				Ticker:     "HIGH",
				Action:     "upgraded by",
				RatingFrom: "Hold",
				RatingTo:   "Strong-Buy",
				TargetFrom: 100.0,
				TargetTo:   200.0,
			},
		},
		{
			name: "Puntuación Media",
			stock: domain.Stock{
				Ticker:     "MED",
				Action:     "target raised by",
				RatingFrom: "Buy",
				RatingTo:   "Buy",
				TargetFrom: 100.0,
				TargetTo:   130.0,
			},
		},
		{
			name: "Puntuación Baja",
			stock: domain.Stock{
				Ticker:     "LOW",
				Action:     "reiterated by",
				RatingFrom: "Hold",
				RatingTo:   "Hold",
				TargetFrom: 100.0,
				TargetTo:   105.0,
			},
		},
	}

	// Verificar que las puntuaciones se ordenan como se espera (mayor a menor)
	highScore := recommendationScore(testCases[0].stock)
	midScore := recommendationScore(testCases[1].stock)
	lowScore := recommendationScore(testCases[2].stock)

	assert.Greater(t, highScore, midScore,
		"La puntuación del caso 'Alta Puntuación' debe ser mayor que 'Puntuación Media'")
	assert.Greater(t, midScore, lowScore,
		"La puntuación del caso 'Puntuación Media' debe ser mayor que 'Puntuación Baja'")
}

// TestNegativeScenarios verifica que el algoritmo maneje correctamente escenarios negativos
func TestNegativeScenarios(t *testing.T) {
	testCases := []struct {
		name  string
		stock domain.Stock
	}{
		{
			name: "Reducción de Precio Objetivo",
			stock: domain.Stock{
				Ticker:     "DOWN",
				Action:     "target lowered by",
				RatingFrom: "Buy",
				RatingTo:   "Buy",
				TargetFrom: 100.0,
				TargetTo:   80.0, // Reducción del 20%
			},
		},
		{
			name: "Degradación de Calificación",
			stock: domain.Stock{
				Ticker:     "DOWNGRADE",
				Action:     "downgraded by",
				RatingFrom: "Buy",
				RatingTo:   "Sell",
				TargetFrom: 100.0,
				TargetTo:   100.0, // Sin cambio en el precio
			},
		},
	}

	// Verificar que las puntuaciones para escenarios negativos son menores que para una acción neutral
	neutralStock := domain.Stock{
		Ticker:     "NEUTRAL",
		Action:     "reiterated by",
		RatingFrom: "Hold",
		RatingTo:   "Hold",
		TargetFrom: 100.0,
		TargetTo:   100.0,
	}
	neutralScore := recommendationScore(neutralStock)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			negativeScore := recommendationScore(tc.stock)
			assert.Less(t, negativeScore, neutralScore,
				"La puntuación para %s debe ser menor que para una acción neutral", tc.name)
		})
	}
}

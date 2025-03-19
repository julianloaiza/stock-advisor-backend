package stocks

import (
	"strings"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// Mapas de puntuaciones preconfigurados
var (
	// Puntuaciones de calificaciones (incluyen valores negativos)
	ratingScores = map[string]float64{
		// Calificaciones muy positivas
		"strong-buy": 25,
		"strong buy": 25,
		"buy":        20,
		"outperform": 18,
		"overweight": 15,

		// Calificaciones positivas moderadas
		"accumulate":        12,
		"add":               12,
		"sector outperform": 10,

		// Calificaciones neutrales
		"market perform": 5,
		"sector perform": 5,
		"equal weight":   5,
		"in-line":        5,
		"hold":           0,
		"neutral":        0,

		// Calificaciones negativas
		"sector weight": -5,
		"market weight": -5,
		"underperform":  -10,
		"underweight":   -10,
		"reduce":        -15,
		"sell":          -20,
		"strong sell":   -25,
	}

	// Fragmentos de texto para acciones y sus puntuaciones
	actionScorePatterns = map[string]float64{
		"upgraded by":       15,  // Mejora en la calificación
		"target raised by":  12,  // Aumento del precio objetivo
		"initiated by":      8,   // Nueva cobertura
		"reiterated by":     5,   // Reiteración
		"target set by":     3,   // Establecimiento de precio objetivo
		"target lowered by": -10, // Reducción del precio objetivo
		"downgraded by":     -12, // Degradación en la calificación
	}

	// Calificaciones consideradas negativas para ajustes posteriores
	negativeRatings = map[string]bool{
		"underperform": true,
		"underweight":  true,
		"reduce":       true,
		"sell":         true,
		"strong sell":  true,
	}
)

// Factores de ponderación para el cálculo de recomendación
const (
	percentDiffWeight   = 0.35 // Cambio porcentual en precio objetivo
	ratingWeight        = 0.30 // Calificación del analista
	actionWeight        = 0.20 // Tipo de acción tomada
	absoluteBonusWeight = 0.15 // Magnitud absoluta del cambio

	// Factores de ajuste
	decreasingTargetFactor = 0.4 // Factor para precios objetivo decrecientes
	negativeRatingFactor   = 0.6 // Factor para calificaciones negativas con score positivo
)

// Estructura para mantener los puntajes base
type baseScoreComponents struct {
	percentDiff   float64
	ratingScore   float64
	actionScore   float64
	absoluteBonus float64
}

// recommendationScore calcula un puntaje de recomendación para una acción.
// Utiliza un enfoque balanceado para evaluar el potencial de inversión.
func (s *service) recommendationScore(stock domain.Stock) float64 {
	// 1. Calcular componentes individuales
	baseScores := s.calculateBaseScores(stock)

	// 2. Aplicar factores externos (empresas, brokerages, etc.)
	adjustedScores := s.applyExternalFactors(stock, baseScores)

	// 3. Calcular puntuación ponderada
	weightedScore := s.calculateWeightedScore(adjustedScores)

	// 4. Aplicar modificadores basados en contexto
	finalScore := s.applyContextModifiers(stock, weightedScore)

	return finalScore
}

// calculateBaseScores calcula los componentes individuales del puntaje
func (s *service) calculateBaseScores(stock domain.Stock) baseScoreComponents {
	return baseScoreComponents{
		percentDiff:   calculatePercentDiff(stock.TargetFrom, stock.TargetTo),
		ratingScore:   calculateRatingScore(stock.RatingTo),
		actionScore:   calculateActionScore(stock.Action),
		absoluteBonus: calculateAbsoluteBonus(stock.TargetFrom, stock.TargetTo),
	}
}

// applyExternalFactors aplica factores externos que pueden afectar los componentes del puntaje
func (s *service) applyExternalFactors(stock domain.Stock, scores baseScoreComponents) baseScoreComponents {
	// Copia los puntajes para no modificar los originales
	adjusted := scores

	// Solo aplicar factores si están disponibles en la configuración
	if s.cfg.RecommendationFactors != nil {
		// Aplicar factor de empresa si existe para este ticker
		if factor, exists := s.cfg.RecommendationFactors.Companies[stock.Ticker]; exists {
			adjusted.percentDiff = scores.percentDiff * (1 + (factor / 100))
		}

		// Aplicar factor de brokerage si existe para este brokerage
		if factor, exists := s.cfg.RecommendationFactors.Brokerages[stock.Brokerage]; exists {
			adjusted.ratingScore = scores.ratingScore * (1 + (factor / 100))
		}
	}

	return adjusted
}

// calculateWeightedScore calcula el puntaje ponderado basado en los componentes
func (s *service) calculateWeightedScore(scores baseScoreComponents) float64 {
	return (scores.percentDiff * percentDiffWeight) +
		(scores.ratingScore * ratingWeight) +
		(scores.actionScore * actionWeight) +
		(scores.absoluteBonus * absoluteBonusWeight)
}

// applyContextModifiers ajusta la puntuación basado en el contexto específico del stock
func (s *service) applyContextModifiers(stock domain.Stock, score float64) float64 {
	adjustedScore := score

	// Reducir puntuación si el precio objetivo está disminuyendo
	if stock.TargetTo < stock.TargetFrom {
		adjustedScore = adjustedScore * decreasingTargetFactor
	}

	// Reducir puntuación para calificaciones negativas cuando el score es positivo
	// Nota: Esto es necesario porque pueden existir casos donde otros factores
	// compensan la calificación negativa, resultando en un score global positivo
	if isNegativeRating(stock.RatingTo) && adjustedScore > 0 {
		adjustedScore = adjustedScore * negativeRatingFactor
	}

	return adjustedScore
}

// isNegativeRating determina si una calificación es considerada negativa
// Usado para aplicar ajustes adicionales a stocks con calificaciones negativas
func isNegativeRating(rating string) bool {
	normalizedRating := strings.ToLower(strings.TrimSpace(rating))
	return negativeRatings[normalizedRating]
}

// calculatePercentDiff calcula la diferencia porcentual entre los precios objetivo
func calculatePercentDiff(from, to float64) float64 {
	if from <= 0 {
		if to > 0 {
			return 20 // Valor positivo moderado si no hay precio inicial
		}
		return 0 // Evitar operaciones con valores inválidos
	}

	// Cálculo simple del cambio porcentual
	percentChange := ((to - from) / from) * 100

	// Penalizar los cambios negativos pero de forma más moderada
	if percentChange < 0 {
		return percentChange * 1.2
	}

	// Valorar más los grandes cambios porcentuales con una escala más gradual
	if percentChange > 50 {
		return 50 + ((percentChange - 50) / 4)
	}

	return percentChange
}

// calculateAbsoluteBonus calcula una bonificación basada en la magnitud absoluta y relativa del cambio
func calculateAbsoluteBonus(from, to float64) float64 {
	diff := to - from

	// Si no tenemos un precio inicial válido, calculamos solo según el precio final
	if from <= 0 {
		if to >= 100 {
			return 15
		} else if to >= 50 {
			return 10
		} else if to >= 20 {
			return 5
		}
		return 0
	}

	// Para cambios negativos, usar una escala más equilibrada
	if diff < 0 {
		absDiff := -diff // Valor absoluto
		relativeDiff := absDiff / from

		if relativeDiff >= 0.20 && from >= 100 { // Caída del 20%+ en stocks de alto valor
			return -15
		} else if relativeDiff >= 0.20 { // Caída del 20%+ en general
			return -12
		} else if relativeDiff >= 0.10 { // Caída del 10%+
			return -8
		} else {
			return -5
		}
	}

	// Para cambios positivos, bonificar más los cambios grandes
	if diff >= 100 {
		return 25
	} else if diff >= 50 {
		return 20
	} else if diff >= 20 {
		return 15
	} else if diff >= 10 {
		return 10
	} else if diff >= 5 {
		return 7
	} else if diff > 0 {
		// Para cambios pequeños pero significativos (ej. de 5 a 7)
		if from <= 10 && diff/from >= 0.15 {
			return 5
		}
		return 3
	}

	return 0
}

// calculateRatingScore evalúa la calificación final del analista
func calculateRatingScore(rating string) float64 {
	normalizedRating := strings.ToLower(strings.TrimSpace(rating))

	// Buscar en el mapa de puntuaciones
	if score, exists := ratingScores[normalizedRating]; exists {
		return score
	}

	return 0 // Valor neutral por defecto
}

// calculateActionScore evalúa el tipo de acción realizada por el analista
func calculateActionScore(action string) float64 {
	normalizedAction := strings.ToLower(strings.TrimSpace(action))

	// Buscar en patrones de acción
	for pattern, score := range actionScorePatterns {
		if strings.Contains(normalizedAction, pattern) {
			return score
		}
	}

	return 0 // Valor neutral por defecto
}

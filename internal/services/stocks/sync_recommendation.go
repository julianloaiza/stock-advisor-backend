package stocks

import (
	"strings"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// recommendationScore calcula un puntaje de recomendación para una acción.
// Utiliza un enfoque balanceado para evaluar el potencial de inversión.
func recommendationScore(s domain.Stock) float64 {
	// Factores de ponderación
	const (
		percentDiffWeight   = 0.35 // Cambio porcentual en precio objetivo
		ratingWeight        = 0.30 // Calificación del analista
		actionWeight        = 0.20 // Tipo de acción tomada
		absoluteBonusWeight = 0.15 // Magnitud absoluta del cambio (aumentado para dar más peso)
	)

	// Calcular componentes individuales
	percentDiff := calculatePercentDiff(s.TargetFrom, s.TargetTo)
	ratingScore := calculateRatingScore(s.RatingTo) // Solo usamos la calificación final
	actionScore := calculateActionScore(s.Action)
	absoluteBonus := calculateAbsoluteBonus(s.TargetFrom, s.TargetTo)

	// Aplicar ponderaciones
	weightedScore := (percentDiff * percentDiffWeight) +
		(ratingScore * ratingWeight) +
		(actionScore * actionWeight) +
		(absoluteBonus * absoluteBonusWeight)

	// Ajuste moderado para precios objetivo decrecientes, en lugar de invertir completamente
	if s.TargetTo < s.TargetFrom {
		// Reducir el puntaje pero no hacer que sea necesariamente negativo
		// Esto permite que otros factores positivos (como una buena calificación)
		// aún tengan influencia
		weightedScore = weightedScore * 0.4
	}

	// Ajuste moderado para calificaciones negativas
	if isNegativeRating(s.RatingTo) && weightedScore > 0 {
		// Reducir el puntaje pero permitir que sea positivo en algunos casos
		weightedScore = weightedScore * 0.6
	}

	return weightedScore
}

// isNegativeRating determina si una calificación es considerada negativa
func isNegativeRating(rating string) bool {
	normalizedRating := strings.ToLower(strings.TrimSpace(rating))

	return normalizedRating == "underperform" ||
		normalizedRating == "underweight" ||
		normalizedRating == "reduce" ||
		normalizedRating == "sell" ||
		normalizedRating == "strong sell"
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
// Esto refleja que un cambio de 100 a 200 es mejor que uno de 10 a 20, aunque porcentualmente sean iguales
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

		// Considerar también el precio base para contextualizar el cambio
		// Un cambio negativo grande en un stock caro es más significativo
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
	// Considerando tanto el valor absoluto como el contexto del precio
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

	// Escala balanceada de calificaciones
	switch normalizedRating {
	// Calificaciones muy positivas
	case "strong-buy", "strong buy":
		return 25
	case "buy":
		return 20
	case "outperform":
		return 18
	case "overweight":
		return 15

	// Calificaciones positivas moderadas
	case "accumulate", "add":
		return 12
	case "sector outperform":
		return 10

	// Calificaciones neutrales
	case "market perform", "sector perform", "equal weight", "in-line":
		return 5
	case "hold", "neutral":
		return 0

	// Calificaciones negativas (valores menos extremos)
	case "sector weight", "market weight":
		return -5
	case "underperform", "underweight":
		return -10
	case "reduce":
		return -15
	case "sell":
		return -20
	case "strong sell":
		return -25

	// Valor por defecto
	default:
		return 0 // Valor neutral por defecto
	}
}

// calculateActionScore evalúa el tipo de acción realizada por el analista
func calculateActionScore(action string) float64 {
	normalizedAction := strings.ToLower(strings.TrimSpace(action))

	// Escala balanceada de acciones
	if strings.Contains(normalizedAction, "upgraded by") {
		return 15 // Mejora en la calificación
	}
	if strings.Contains(normalizedAction, "target raised by") {
		return 12 // Aumento del precio objetivo
	}
	if strings.Contains(normalizedAction, "initiated by") {
		return 8 // Nueva cobertura
	}
	if strings.Contains(normalizedAction, "reiterated by") {
		return 5 // Reiteración
	}
	if strings.Contains(normalizedAction, "target set by") {
		return 3 // Establecimiento de precio objetivo
	}
	if strings.Contains(normalizedAction, "target lowered by") {
		return -10 // Reducción del precio objetivo (valor menos extremo)
	}
	if strings.Contains(normalizedAction, "downgraded by") {
		return -12 // Degradación en la calificación (valor menos extremo)
	}

	return 0
}

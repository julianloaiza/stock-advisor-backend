package stocks

import (
	"strings"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// recommendationScore calcula un puntaje de recomendación para una acción.
func recommendationScore(s domain.Stock) float64 {
	// Calcular la diferencia porcentual entre target_to y target_from
	percentDiff := calculatePercentDiff(s.TargetFrom, s.TargetTo)

	// Calcular bonificaciones adicionales
	absoluteBonus := calculateAbsoluteBonus(s.TargetTo - s.TargetFrom)
	ratingBonus := calculateRatingBonus(s.RatingTo)
	actionBonus := calculateActionBonus(s.Action)

	// Retornar la puntuación total
	return percentDiff + absoluteBonus + ratingBonus + actionBonus
}

// calculatePercentDiff calcula la diferencia porcentual
func calculatePercentDiff(from, to float64) float64 {
	if from != 0 {
		return ((to - from) / from) * 100
	}
	return to // Si from es 0, usamos to como diferencia
}

// calculateAbsoluteBonus calcula la bonificación basada en la diferencia absoluta
func calculateAbsoluteBonus(diff float64) float64 {
	switch {
	case diff >= 100:
		return 10
	case diff >= 50:
		return 7
	case diff >= 20:
		return 5
	case diff >= 10:
		return 3
	case diff >= 5:
		return 2
	default:
		return 0
	}
}

// calculateRatingBonus calcula la bonificación basada en la calificación
func calculateRatingBonus(rating string) float64 {
	switch strings.ToLower(rating) {
	case "buy":
		return 10
	case "strong-buy":
		return 15
	case "outperform":
		return 12
	default:
		return 0
	}
}

// calculateActionBonus calcula la bonificación basada en la acción
func calculateActionBonus(action string) float64 {
	switch strings.ToLower(action) {
	case "target raised by":
		return 5
	case "upgraded by":
		return 7
	default:
		return 0
	}
}

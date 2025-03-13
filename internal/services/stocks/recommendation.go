package stocks

import (
	"strings"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// recommendationScore calcula un puntaje de recomendación para una acción.
func recommendationScore(s domain.Stock) float64 {
	var percentDiff float64
	if s.TargetFrom != 0 {
		percentDiff = ((s.TargetTo - s.TargetFrom) / s.TargetFrom) * 100
	} else {
		percentDiff = s.TargetTo
	}

	absoluteDiff := s.TargetTo - s.TargetFrom

	var absoluteBonus float64
	switch {
	case absoluteDiff >= 100:
		absoluteBonus = 10
	case absoluteDiff >= 50:
		absoluteBonus = 7
	case absoluteDiff >= 20:
		absoluteBonus = 5
	case absoluteDiff >= 10:
		absoluteBonus = 3
	case absoluteDiff >= 5:
		absoluteBonus = 2
	default:
		absoluteBonus = 0
	}

	var ratingBonus float64
	switch strings.ToLower(s.RatingTo) {
	case "buy":
		ratingBonus = 10
	case "strong-buy":
		ratingBonus = 15
	case "outperform":
		ratingBonus = 12
	default:
		ratingBonus = 0
	}

	var actionBonus float64
	switch strings.ToLower(s.Action) {
	case "target raised by":
		actionBonus = 5
	case "upgraded by":
		actionBonus = 7
	default:
		actionBonus = 0
	}

	return percentDiff + absoluteBonus + ratingBonus + actionBonus
}

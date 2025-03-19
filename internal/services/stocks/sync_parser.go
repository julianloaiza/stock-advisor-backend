package stocks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// parseStock convierte un mapa a un objeto domain.Stock y le asigna una puntuación de recomendación
func (s *service) parseStock(item map[string]interface{}) (domain.Stock, error) {
	// Extraer campos de texto del mapa
	textFields := s.extractTextFields(item)

	// Procesar campos numéricos
	targetFrom, targetTo, err := s.extractNumericFields(item)
	if err != nil {
		return domain.Stock{}, err
	}

	// Crear objeto Stock
	stock := domain.Stock{
		Ticker:     textFields["ticker"],
		Company:    textFields["company"],
		Brokerage:  textFields["brokerage"],
		Action:     textFields["action"],
		RatingFrom: textFields["ratingFrom"],
		RatingTo:   textFields["ratingTo"],
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Currency:   textFields["currency"],
	}

	// Calcular y asignar la puntuación de recomendación
	stock.RecommendScore = s.recommendationScore(stock)

	return stock, nil
}

// extractTextFields extrae los campos de texto del mapa de datos
func (s *service) extractTextFields(item map[string]interface{}) map[string]string {
	fields := make(map[string]string)

	// Extraer valores de texto del mapa
	fields["ticker"], _ = item["ticker"].(string)
	fields["company"], _ = item["company"].(string)
	fields["brokerage"], _ = item["brokerage"].(string)
	fields["action"], _ = item["action"].(string)
	fields["ratingFrom"], _ = item["rating_from"].(string)
	fields["ratingTo"], _ = item["rating_to"].(string)
	fields["currency"], _ = item["currency"].(string)

	// Valores por defecto
	if fields["currency"] == "" {
		fields["currency"] = "USD"
	}

	return fields
}

// extractNumericFields procesa y extrae los campos numéricos
func (s *service) extractNumericFields(item map[string]interface{}) (float64, float64, error) {
	// Procesar valores numéricos
	targetFromStr, _ := item["target_from"].(string)
	targetToStr, _ := item["target_to"].(string)

	// Limpiar formatos monetarios
	targetFromStr = s.cleanMonetaryFormat(targetFromStr)
	targetToStr = s.cleanMonetaryFormat(targetToStr)

	// Convertir target_from a número
	targetFrom, err := strconv.ParseFloat(targetFromStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error convirtiendo target_from: %w", err)
	}

	// Convertir target_to a número
	targetTo, err := strconv.ParseFloat(targetToStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error convirtiendo target_to: %w", err)
	}

	return targetFrom, targetTo, nil
}

// cleanMonetaryFormat elimina símbolos de moneda y separadores de miles
func (s *service) cleanMonetaryFormat(value string) string {
	return strings.ReplaceAll(strings.TrimPrefix(value, "$"), ",", "")
}

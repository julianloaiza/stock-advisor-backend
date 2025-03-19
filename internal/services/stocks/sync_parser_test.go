package stocks

import (
	"testing"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/stretchr/testify/assert"
)

// TestExtractTextFields verifica la extracción de campos de texto
func TestExtractTextFields(t *testing.T) {
	// Crear una configuración básica para el servicio
	cfg := &config.Config{
		SyncMaxIterations: 10,
		SyncTimeout:       60,
	}

	// Crear instancia del servicio para pruebas
	s := &service{
		cfg: cfg,
	}

	// Caso completo
	item := map[string]interface{}{
		"ticker":      "AAPL",
		"company":     "Apple Inc.",
		"brokerage":   "Example Broker",
		"action":      "upgraded by",
		"rating_from": "Hold",
		"rating_to":   "Buy",
		"currency":    "EUR",
	}

	fields := s.extractTextFields(item)

	assert.Equal(t, "AAPL", fields["ticker"])
	assert.Equal(t, "Apple Inc.", fields["company"])
	assert.Equal(t, "Example Broker", fields["brokerage"])
	assert.Equal(t, "upgraded by", fields["action"])
	assert.Equal(t, "Hold", fields["ratingFrom"])
	assert.Equal(t, "Buy", fields["ratingTo"])
	assert.Equal(t, "EUR", fields["currency"])

	// Caso con campos faltantes
	itemIncomplete := map[string]interface{}{
		"ticker":  "MSFT",
		"company": "Microsoft",
		// Otros campos faltantes
	}

	fieldsIncomplete := s.extractTextFields(itemIncomplete)

	assert.Equal(t, "MSFT", fieldsIncomplete["ticker"])
	assert.Equal(t, "Microsoft", fieldsIncomplete["company"])
	assert.Equal(t, "", fieldsIncomplete["brokerage"])
	assert.Equal(t, "USD", fieldsIncomplete["currency"], "La moneda debe ser USD por defecto")
}

// TestExtractNumericFields verifica la extracción de campos numéricos
func TestExtractNumericFields(t *testing.T) {
	// Crear una configuración básica para el servicio
	cfg := &config.Config{
		SyncMaxIterations: 10,
		SyncTimeout:       60,
	}

	// Crear instancia del servicio para pruebas
	s := &service{
		cfg: cfg,
	}

	// Caso valores correctos
	t.Run("Valores correctos", func(t *testing.T) {
		item := map[string]interface{}{
			"target_from": "100.50",
			"target_to":   "150.75",
		}

		from, to, err := s.extractNumericFields(item)

		assert.NoError(t, err)
		assert.Equal(t, 100.50, from)
		assert.Equal(t, 150.75, to)
	})

	// Caso formato con símbolo de moneda y separadores
	t.Run("Formato monetario", func(t *testing.T) {
		item := map[string]interface{}{
			"target_from": "$1,234.56",
			"target_to":   "$2,345.67",
		}

		from, to, err := s.extractNumericFields(item)

		assert.NoError(t, err)
		assert.Equal(t, 1234.56, from)
		assert.Equal(t, 2345.67, to)
	})

	// Caso error en target_from
	t.Run("Error en target_from", func(t *testing.T) {
		item := map[string]interface{}{
			"target_from": "invalid",
			"target_to":   "100.00",
		}

		_, _, err := s.extractNumericFields(item)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error convirtiendo target_from")
	})

	// Caso error en target_to
	t.Run("Error en target_to", func(t *testing.T) {
		item := map[string]interface{}{
			"target_from": "100.00",
			"target_to":   "invalid",
		}

		_, _, err := s.extractNumericFields(item)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error convirtiendo target_to")
	})
}

// TestCleanMonetaryFormat verifica la limpieza de formatos monetarios
func TestCleanMonetaryFormat(t *testing.T) {
	// Crear una configuración básica para el servicio
	cfg := &config.Config{
		SyncMaxIterations: 10,
		SyncTimeout:       60,
	}

	// Crear instancia del servicio para pruebas
	s := &service{
		cfg: cfg,
	}

	testCases := []struct {
		input    string
		expected string
	}{
		{"$1,234.56", "1234.56"},
		{"1,234.56", "1234.56"},
		{"$123.45", "123.45"},
		{"123.45", "123.45"},
		{"$1,000", "1000"},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := s.cleanMonetaryFormat(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

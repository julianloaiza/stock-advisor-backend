package stocks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// SyncStocks sincroniza la base de datos con la API externa.
func (s *service) SyncStocks(ctx context.Context, limit int) error {
	// Crear un contexto con timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.SyncTimeout)*time.Second)
	defer cancel()

	// Ejecutar sincronización
	return s.syncStocks(ctx, limit)
}

// syncStocks es la implementación principal de la sincronización
func (s *service) syncStocks(ctx context.Context, limit int) error {
	// Validar y ajustar el límite de iteraciones
	limit = s.validateLimit(limit)
	allStocks := make([]domain.Stock, 0, limit*10)

	log.Println("🔄 Iniciando sincronización con la API")

	// Variables para control de iteración
	var nextPage string
	seenTokens := make(map[string]bool)

	// Iterar para obtener datos paginados
	for i := 1; i <= limit; i++ {
		// Preparar parámetros para la solicitud
		params := make(map[string]string)
		if nextPage != "" {
			params["next_page"] = nextPage
		}

		// Realizar solicitud a la API
		responseData, err := s.apiClient.Get(ctx, "", params)
		if err != nil {
			return fmt.Errorf("error en iteración %d: %w", i, err)
		}

		// Definir estructura para la respuesta
		var result struct {
			Items    []map[string]interface{} `json:"items"`
			NextPage string                   `json:"next_page"`
		}

		// Deserializar la respuesta
		if err := json.Unmarshal(responseData, &result); err != nil {
			return fmt.Errorf("error parseando JSON en iteración %d: %w", i, err)
		}

		log.Printf("Iteración %d: next_page value = %s", i, result.NextPage)

		// Procesar elementos
		for _, item := range result.Items {
			stock, err := s.parseStock(item)
			if err != nil {
				log.Printf("Iteración %d: error parseando stock: %v", i, err)
				continue // Continuar con el siguiente item en caso de error
			}
			allStocks = append(allStocks, stock)
		}

		// Verificar si debemos terminar la sincronización
		if s.shouldTerminateSync(result.NextPage, seenTokens) {
			break
		}

		// Preparar siguiente iteración
		seenTokens[result.NextPage] = true
		nextPage = result.NextPage
	}

	// Guardar en base de datos
	return s.replaceAllStocks(allStocks)
}

// validateLimit valida el parámetro limit y lo ajusta si es necesario.
func (s *service) validateLimit(limit int) int {
	maxIterations := s.cfg.SyncMaxIterations
	if limit <= 0 {
		log.Printf("Límite inválido (%d). Se usará el valor por defecto: 1", limit)
		return 1
	}
	if limit > maxIterations {
		log.Printf("Límite (%d) excede el máximo permitido (%d). Se ajustará a %d",
			limit, maxIterations, maxIterations)
		return maxIterations
	}
	return limit
}

// shouldTerminateSync determina si la sincronización debe terminar
func (s *service) shouldTerminateSync(nextPage string, seenTokens map[string]bool) bool {
	// Sin más páginas
	if nextPage == "" {
		log.Println("No se recibió next_page. Finalizando sincronización.")
		return true
	}

	// Ciclo detectado
	if seenTokens[nextPage] {
		log.Printf("Detectado ciclo: next_page '%s' ya fue visto. Finalizando sincronización.", nextPage)
		return true
	}

	return false
}

// replaceAllStocks reemplaza todos los stocks en la base de datos
func (s *service) replaceAllStocks(allStocks []domain.Stock) error {
	if len(allStocks) == 0 {
		log.Println("No se encontraron stocks para sincronizar.")
		return nil
	}

	if err := s.repo.ReplaceAllStocks(allStocks); err != nil {
		return fmt.Errorf("error reemplazando stocks: %w", err)
	}

	log.Printf("Sincronización completada exitosamente. %d stocks sincronizados.", len(allStocks))
	return nil
}

// parseStock convierte un mapa a un objeto domain.Stock y le asigna una puntuación de recomendación
func (s *service) parseStock(item map[string]interface{}) (domain.Stock, error) {
	var stock domain.Stock

	// Extraer valores del mapa
	ticker, _ := item["ticker"].(string)
	company, _ := item["company"].(string)
	brokerage, _ := item["brokerage"].(string)
	action, _ := item["action"].(string)
	ratingFrom, _ := item["rating_from"].(string)
	ratingTo, _ := item["rating_to"].(string)
	currency, _ := item["currency"].(string)

	// Valores por defecto
	if currency == "" {
		currency = "USD"
	}

	// Procesar valores numéricos
	targetFromStr, _ := item["target_from"].(string)
	targetToStr, _ := item["target_to"].(string)

	// Limpiar formatos monetarios
	targetFromStr = strings.ReplaceAll(strings.TrimPrefix(targetFromStr, "$"), ",", "")
	targetToStr = strings.ReplaceAll(strings.TrimPrefix(targetToStr, "$"), ",", "")

	// Convertir a números
	targetFrom, err := strconv.ParseFloat(targetFromStr, 64)
	if err != nil {
		return stock, fmt.Errorf("error convirtiendo target_from: %w", err)
	}

	targetTo, err := strconv.ParseFloat(targetToStr, 64)
	if err != nil {
		return stock, fmt.Errorf("error convirtiendo target_to: %w", err)
	}

	// Crear objeto Stock
	stock = domain.Stock{
		Ticker:     ticker,
		Company:    company,
		Brokerage:  brokerage,
		Action:     action,
		RatingFrom: ratingFrom,
		RatingTo:   ratingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Currency:   currency,
	}

	// Calcular y asignar la puntuación de recomendación
	stock.RecommendScore = recommendationScore(stock)

	return stock, nil
}

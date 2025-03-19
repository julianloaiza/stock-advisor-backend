package stocks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		// Obtener datos de la página actual
		items, newNextPage, err := s.fetchPageData(ctx, nextPage, i)
		if err != nil {
			return err
		}

		// Procesar elementos
		pageStocks := s.processPageItems(items, i)
		allStocks = append(allStocks, pageStocks...)

		// Verificar si debemos terminar la sincronización
		if s.shouldTerminateSync(newNextPage, seenTokens) {
			break
		}

		// Preparar siguiente iteración
		seenTokens[newNextPage] = true
		nextPage = newNextPage
	}

	// Guardar en base de datos
	return s.replaceAllStocks(allStocks)
}

// fetchPageData obtiene los datos de una página de la API
func (s *service) fetchPageData(ctx context.Context, nextPage string, iteration int) ([]map[string]interface{}, string, error) {
	// Preparar parámetros para la solicitud
	params := make(map[string]string)
	if nextPage != "" {
		params["next_page"] = nextPage
	}

	// Realizar solicitud a la API
	responseData, err := s.apiClient.Get(ctx, "", params)
	if err != nil {
		return nil, "", fmt.Errorf("error en iteración %d: %w", iteration, err)
	}

	// Definir estructura para la respuesta
	var result struct {
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page"`
	}

	// Deserializar la respuesta
	if err := json.Unmarshal(responseData, &result); err != nil {
		return nil, "", fmt.Errorf("error parseando JSON en iteración %d: %w", iteration, err)
	}

	log.Printf("Iteración %d: next_page value = %s", iteration, result.NextPage)
	return result.Items, result.NextPage, nil
}

// processPageItems procesa los elementos de una página y los convierte a stocks
func (s *service) processPageItems(items []map[string]interface{}, iteration int) []domain.Stock {
	var pageStocks []domain.Stock
	for _, item := range items {
		stock, err := s.parseStock(item)
		if err != nil {
			log.Printf("Iteración %d: error parseando stock: %v", iteration, err)
			continue // Continuar con el siguiente item en caso de error
		}
		pageStocks = append(pageStocks, stock)
	}
	return pageStocks
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

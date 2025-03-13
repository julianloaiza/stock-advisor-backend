package stocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// SyncStocks utiliza la función auxiliar (definida en sync.go) para sincronizar la base de datos.
func (s *service) SyncStocks(ctx context.Context, limit int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.SyncTimeout)*time.Second)
	defer cancel()
	return s.syncStocks(ctx, limit)
}

// syncStocks es la función auxiliar que contiene la lógica para sincronizar
// la base de datos con la API externa.
func (s *service) syncStocks(ctx context.Context, limit int) error {
	limit = s.validateLimit(limit)
	allStocks := make([]domain.Stock, 0, limit*10)

	log.Println("🔄 Iniciando sincronización con la API")
	client := &http.Client{}
	baseURL := s.cfg.StockAPIURL
	authToken := "Bearer " + s.cfg.StockAPIKey

	var nextPage string
	seenTokens := make(map[string]bool)

	for i := 1; i <= limit; i++ {
		url := s.constructURL(baseURL, nextPage)
		resp, err := s.makeRequest(ctx, client, url, authToken)
		if err != nil {
			return err
		}

		body, err := s.readResponseBody(resp)
		if err != nil {
			return err
		}

		result, err := s.parseResponseBody(body)
		if err != nil {
			return err
		}

		log.Printf("Iteración %d: next_page value = %s", i, result.NextPage)

		for _, item := range result.Items {
			stock, err := parseStock(item)
			if err != nil {
				log.Printf("Iteración %d: error parseando stock: %v", i, err)
				return err
			}
			allStocks = append(allStocks, stock)
		}

		if s.shouldTerminateSync(result.NextPage, seenTokens) {
			break
		}
		seenTokens[result.NextPage] = true
		nextPage = result.NextPage
	}

	return s.replaceAllStocks(allStocks)
}

// validateLimit valida el parámetro limit y lo ajusta si es necesario.
func (s *service) validateLimit(limit int) int {
	maxIterations := s.cfg.SyncMaxIterations
	if limit > maxIterations {
		log.Printf("El parámetro limit (%d) excede el máximo permitido (%d). Se utilizarán %d iteraciones.",
			limit, maxIterations, maxIterations)
		limit = maxIterations
	}
	return limit
}

// constructURL construye la URL de la API con el parámetro next_page.
func (s *service) constructURL(baseURL, nextPage string) string {
	if nextPage != "" {
		return fmt.Sprintf("%s?next_page=%s", baseURL, nextPage)
	}
	return baseURL
}

// makeRequest realiza una solicitud HTTP GET con el cliente proporcionado.
func (s *service) makeRequest(ctx context.Context, client *http.Client, url, authToken string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creando la request: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error realizando la request: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		errMsg := fmt.Sprintf("Status code inesperado: %d", resp.StatusCode)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	return resp, nil
}

// readResponseBody lee el cuerpo de la respuesta HTTP.
func (s *service) readResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("Error leyendo la respuesta: %v", err)
		return nil, err
	}
	return body, nil
}

// parseResponseBody decodifica el cuerpo de la respuesta JSON.
func (s *service) parseResponseBody(body []byte) (*struct {
	Items    []map[string]interface{} `json:"items"`
	NextPage string                   `json:"next_page"`
}, error) {
	var result struct {
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error decodificando JSON: %v", err)
		return nil, err
	}
	return &result, nil
}

// shouldTerminateSync determina si la sincronización debe terminar.
func (s *service) shouldTerminateSync(nextPage string, seenTokens map[string]bool) bool {
	if nextPage == "" {
		log.Println("No se recibió next_page. Finalizando sincronización.")
		return true
	}
	if seenTokens[nextPage] {
		log.Printf("Detectado ciclo: next_page '%s' ya fue visto. Finalizando sincronización.", nextPage)
		return true
	}
	return false
}

// replaceAllStocks reemplaza todos los stocks en la base de datos.
func (s *service) replaceAllStocks(allStocks []domain.Stock) error {
	if err := s.repo.ReplaceAllStocks(allStocks); err != nil {
		log.Printf("Error reemplazando stocks: %v", err)
		return err
	}
	log.Println("Sincronización completada exitosamente.")
	return nil
}

// parseStock convierte un mapa (map[string]interface{}) en un objeto domain.Stock.
func parseStock(item map[string]interface{}) (domain.Stock, error) {
	var s domain.Stock

	ticker, _ := item["ticker"].(string)
	company, _ := item["company"].(string)
	brokerage, _ := item["brokerage"].(string)
	action, _ := item["action"].(string)
	ratingFrom, _ := item["rating_from"].(string)
	ratingTo, _ := item["rating_to"].(string)
	timeStr, _ := item["time"].(string)
	targetFromStr, _ := item["target_from"].(string)
	targetToStr, _ := item["target_to"].(string)

	targetFromStr = strings.ReplaceAll(strings.TrimPrefix(targetFromStr, "$"), ",", "")
	targetToStr = strings.ReplaceAll(strings.TrimPrefix(targetToStr, "$"), ",", "")

	targetFrom, err := strconv.ParseFloat(targetFromStr, 64)
	if err != nil {
		return s, fmt.Errorf("error converting target_from: %v", err)
	}
	targetTo, err := strconv.ParseFloat(targetToStr, 64)
	if err != nil {
		return s, fmt.Errorf("error converting target_to: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return s, fmt.Errorf("error parsing time: %v", err)
	}

	s = domain.Stock{
		Ticker:     ticker,
		Company:    company,
		Brokerage:  brokerage,
		Action:     action,
		RatingFrom: ratingFrom,
		RatingTo:   ratingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Time:       parsedTime,
		Currency:   "USD",
	}
	return s, nil
}

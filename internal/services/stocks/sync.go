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

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/internal/domain"
	repo "github.com/julianloaiza/stock-advisor/internal/repositories/stocks"
)

// syncStocks es la función auxiliar que contiene la lógica para sincronizar
// la base de datos con la API externa. Acumula en memoria todos los registros
// obtenidos y, al finalizar, reemplaza la data antigua en una única operación.
//
// Se valida que el parámetro "limit" no exceda el máximo permitido (cfg.SyncMaxIterations)
// y se preasigna el slice con capacidad = limit * 10.
func syncStocks(ctx context.Context, limit int, repository repo.Repository, cfg *config.Config) error {
	// Validar que "limit" no exceda el máximo permitido.
	maxIterations := cfg.SyncMaxIterations // Por ejemplo, 100 (definido en .env como SYNC_MAX_ITERATIONS)
	if limit > maxIterations {
		log.Printf("El parámetro limit (%d) excede el máximo permitido (%d). Se utilizarán %d iteraciones.",
			limit, maxIterations, maxIterations)
		limit = maxIterations
	}

	// Preasignar el slice con capacidad = limit * 10.
	allStocks := make([]domain.Stock, 0, limit*10)

	log.Println("🔄 Iniciando sincronización con la API")
	client := &http.Client{}
	baseURL := cfg.StockAPIURL
	authToken := "Bearer " + cfg.StockAPIKey

	var nextPage string
	seenTokens := make(map[string]bool)

	// Iterar hasta el límite definido.
	for i := 1; i <= limit; i++ {
		// Construir la URL con el query param "next_page" si corresponde.
		url := baseURL
		if nextPage != "" {
			url = fmt.Sprintf("%s?next_page=%s", baseURL, nextPage)
		}

		// Crear la request HTTP con contexto.
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Printf("Iteración %d: error creando la request: %v", i, err)
			return err
		}
		req.Header.Set("Authorization", authToken)
		req.Header.Set("Content-Type", "application/json")

		// Ejecutar la request.
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Iteración %d: error realizando la request: %v", i, err)
			return err
		}

		// Verificar el status code.
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			errMsg := fmt.Sprintf("Iteración %d: status code inesperado: %d", i, resp.StatusCode)
			log.Println(errMsg)
			return fmt.Errorf(errMsg)
		}

		// Leer y cerrar el body.
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Iteración %d: error leyendo la respuesta: %v", i, err)
			return err
		}

		// Estructura para decodificar la respuesta.
		var result struct {
			Items    []map[string]interface{} `json:"items"`
			NextPage string                   `json:"next_page"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Printf("Iteración %d: error decodificando JSON: %v", i, err)
			return err
		}

		log.Printf("Iteración %d: next_page value = %s", i, result.NextPage)

		// Convertir cada item a domain.Stock y acumular.
		for _, item := range result.Items {
			stock, err := parseStock(item)
			if err != nil {
				log.Printf("Iteración %d: error parseando stock: %v", i, err)
				return err
			}
			allStocks = append(allStocks, stock)
		}

		// Si no se recibe next_page o se detecta ciclo, finalizar.
		if result.NextPage == "" {
			log.Println("No se recibió next_page. Finalizando sincronización.")
			break
		}
		if seenTokens[result.NextPage] {
			log.Printf("Detectado ciclo en iteración %d: next_page '%s' ya fue visto. Finalizando sincronización.", i, result.NextPage)
			break
		}
		seenTokens[result.NextPage] = true
		nextPage = result.NextPage
	}

	// Reemplazar toda la data en la base de datos en una única transacción.
	if err := repository.ReplaceAllStocks(allStocks); err != nil {
		log.Printf("Error reemplazando stocks: %v", err)
		return err
	}

	log.Println("Sincronización completada exitosamente.")
	return nil
}

// parseStock convierte un mapa (map[string]interface{}) en un objeto domain.Stock.
// Realiza las conversiones necesarias para los campos numéricos y de fecha.
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

	// Eliminar el símbolo "$" y las comas de los valores.
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

package stocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// syncStocks implementa la lógica para sincronizar la base de datos con la API.
func syncStocks(ctx context.Context) error {
	log.Println("🔄 Iniciando sincronización con la API de Truora")

	client := &http.Client{}
	baseURL := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	authToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MSwiZW1haWwiOiJsb2FpemFqdWxpYW4xOTk5QGdtYWlsLmNvbSIsImV4cCI6MTc0MTQ1Mjc0OCwiaWQiOiIwIiwicGFzc3dvcmQiOiInIE9SICcxJz0nMSJ9.adVaiW9LmcuxjPC4kclyMB7bjUZVKbJxmVj1qLobtLI"

	var nextPage string
	// Mapa para almacenar los tokens ya vistos y detectar ciclos.
	seenTokens := make(map[string]bool)

	// Iterar hasta 100 veces.
	for i := 1; i <= 10000; i++ {
		// Construir la URL con el query param si ya tenemos un next_page
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

		// Realizar la request.
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Iteración %d: error realizando la request: %v", i, err)
			return err
		}

		// Leer y decodificar el body de la respuesta.
		body, err := ioutil.ReadAll(resp.Body)
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

		// Imprimir el valor de next_page de esta iteración.
		log.Printf("Iteración %d: next_page value = %s", i, result.NextPage)

		// Validar: si next_page está vacío, finalizamos.
		if result.NextPage == "" {
			log.Println("No se recibió next_page. Finalizando sincronización.")
			break
		}

		// Validar si ya se vio este token (detecta ciclo).
		if seenTokens[result.NextPage] {
			log.Printf("Detectado ciclo en iteración %d: next_page '%s' ya fue visto. Finalizando sincronización.", i, result.NextPage)
			break
		}

		// Almacenar el token para evitar ciclos.
		seenTokens[result.NextPage] = true
		// Actualizar nextPage para la siguiente iteración.
		nextPage = result.NextPage
	}

	log.Println("Sincronización finalizada.")
	return nil
}

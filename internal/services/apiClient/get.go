package apiClient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Get implementa la interfaz Client.Get
func (c *client) Get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	// Construir URL completa
	endpoint := c.buildURL(path, params)

	// Crear la solicitud
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request GET: %w", err)
	}

	// Añadir headers
	c.addHeaders(req)

	// Ejecutar solicitud
	return c.executeRequest(req)
}

// buildURL construye la URL completa con parámetros de consulta
func (c *client) buildURL(path string, params map[string]string) string {
	baseURL := c.baseURL
	if path != "" {
		baseURL = fmt.Sprintf("%s/%s", baseURL, path)
	}

	// Si no hay parámetros, devolver la URL base
	if len(params) == 0 {
		return baseURL
	}

	// Construir la cadena de consulta
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return fmt.Sprintf("%s?%s", baseURL, values.Encode())
}

// addHeaders añade los headers de autenticación y otros comunes
func (c *client) addHeaders(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set(c.authHeader, c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

// executeRequest ejecuta la solicitud HTTP y maneja la respuesta
func (c *client) executeRequest(req *http.Request) ([]byte, error) {
	// Ejecutar solicitud
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de respuesta
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code inesperado: %d, respuesta: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Leer la respuesta
	return io.ReadAll(resp.Body)
}

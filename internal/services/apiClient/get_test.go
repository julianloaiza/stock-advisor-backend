package apiClient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/stretchr/testify/assert"
)

func TestGet_SuccessResponse(t *testing.T) {
	// Crear un servidor de prueba que devuelve una respuesta exitosa
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar método
		assert.Equal(t, http.MethodGet, r.Method)

		// Verificar headers
		assert.Equal(t, "test-auth-tkn", r.Header.Get("X-AUTH-TKN"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verificar parámetros de consulta
		assert.Equal(t, "value1", r.URL.Query().Get("param1"))
		assert.Equal(t, "value2", r.URL.Query().Get("param2"))

		// Responder con JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"test"}]}`))
	}))
	defer server.Close()

	// Crear config con la URL del servidor de prueba
	cfg := &config.Config{
		StockAPIURL:  server.URL,
		StockAuthTkn: "test-auth-tkn",
	}

	// Crear cliente API
	client := &client{
		httpClient: server.Client(),
		baseURL:    cfg.StockAPIURL,
		authToken:  cfg.StockAuthTkn,
	}

	// Ejecutar solicitud GET con parámetros
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	response, err := client.Get(context.Background(), "", params)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Contains(t, string(response), `"success":true`)
}

func TestGet_ErrorResponse(t *testing.T) {
	// Crear un servidor de prueba que devuelve un error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	// Crear config con la URL del servidor de prueba
	cfg := &config.Config{
		StockAPIURL:  server.URL,
		StockAuthTkn: "test-auth-tkn",
	}

	// Crear cliente API
	client := &client{
		httpClient: server.Client(),
		baseURL:    cfg.StockAPIURL,
		authToken:  cfg.StockAuthTkn,
	}

	// Ejecutar solicitud GET
	response, err := client.Get(context.Background(), "", nil)

	// Verificar que se devuelva un error
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "status code inesperado: 500")
}

func TestBuildURL_WithParams(t *testing.T) {
	// Crear cliente API
	client := &client{
		baseURL: "https://api.example.com",
	}

	// Caso 1: URL base sin path con parámetros
	params1 := map[string]string{
		"page": "1",
		"size": "10",
	}
	url1 := client.buildURL("", params1)
	assert.Equal(t, "https://api.example.com?page=1&size=10", url1)

	// Caso 2: URL con path y parámetros
	params2 := map[string]string{
		"query": "test",
	}
	url2 := client.buildURL("search", params2)
	assert.Equal(t, "https://api.example.com/search?query=test", url2)

	// Caso 3: URL sin parámetros
	url3 := client.buildURL("items", nil)
	assert.Equal(t, "https://api.example.com/items", url3)
}

func TestAddHeaders(t *testing.T) {
	// Crear cliente API con API key
	apiClient := &client{
		authToken: "test-auth-tkn",
	}

	// Crear request
	req, _ := http.NewRequest(http.MethodGet, "https://api.example.com", nil)

	// Añadir headers
	apiClient.addHeaders(req)

	// Verificar headers
	assert.Equal(t, "test-auth-tkn", req.Header.Get("X-AUTH-TKN"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, "application/json", req.Header.Get("Accept"))

	// Probar sin API key
	apiClientNoKey := &client{
		authToken: "",
	}

	req2, _ := http.NewRequest(http.MethodGet, "https://api.example.com", nil)
	apiClientNoKey.addHeaders(req2)

	// Verificar que no se agregó el header de autenticación
	assert.Equal(t, "", req2.Header.Get("X-AUTH-TKN"))
}

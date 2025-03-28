package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Crear una instancia de Echo
	e := echo.New()

	// Crear la solicitud GET a /health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Instanciar el handler de health
	h := &handler{}

	// Ejecutar el handler directamente
	err := h.HealthCheck(c)
	assert.NoError(t, err)

	// Verificar que se retorne el status 200
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verificar que el body tenga el JSON esperado
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

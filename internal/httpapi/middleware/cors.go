package middleware

import (
	"strings"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ApplyCORS configura y aplica CORS en Echo basado en .env
func ApplyCORS(e *echo.Echo, cfg *config.Config) {
	allowedOrigins := strings.Split(cfg.CORSAllowedOrigins, ",") // ✅ Soporte para múltiples orígenes separados por ","

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}

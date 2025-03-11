package handlers

import "github.com/labstack/echo/v4"

// Handler es la interfaz que deben implementar todos los controladores HTTP.
type Handler interface {
	RegisterRoutes(e *echo.Echo)
}

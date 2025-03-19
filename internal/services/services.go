package services

import (
	"github.com/julianloaiza/stock-advisor/internal/services/apiClient"
	"github.com/julianloaiza/stock-advisor/internal/services/stocks"
	"go.uber.org/fx"
)

// Module registra los servicios.
var Module = fx.Module("services", fx.Provide(
	apiClient.New, // Servicio API para comunicación con servicios externos
	stocks.New,    // Servicio de stocks
))

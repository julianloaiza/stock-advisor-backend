package httpapi

import (
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/stocks"
	"go.uber.org/fx"
)

// Module registra los handlers de la API.
var Module = fx.Module("httpapi", fx.Provide(
	stocks.New,
))

package services

import (
	"github.com/julianloaiza/stock-advisor/internal/services/stocks"
	"go.uber.org/fx"
)

// Module registra los servicios.
var Module = fx.Module("services", fx.Provide(
	stocks.New,
))

package repositories

import (
	"github.com/julianloaiza/stock-advisor/internal/repositories/stocks"
	"go.uber.org/fx"
)

// Module registra los repositorios.
var Module = fx.Module("repositories", fx.Provide(
	stocks.New,
))

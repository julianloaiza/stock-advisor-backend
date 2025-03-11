package repositories

import (
	"github.com/julianloaiza/stock-advisor/internal/repositories/stocks"
	"go.uber.org/fx"
)

var Module = fx.Module("repositories", fx.Provide(
	stocks.New,
))

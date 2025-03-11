package services

import (
	"github.com/julianloaiza/stock-advisor/internal/services/stocks"
	"go.uber.org/fx"
)

var Module = fx.Module("services", fx.Provide(
	stocks.New,
))

package services

import (
	"github.com/julianloaiza/stock-advisor/internal/services/auth"
	"github.com/julianloaiza/stock-advisor/internal/services/profiles"
	"go.uber.org/fx"
)

var Module = fx.Module("services", fx.Provide(
	auth.New,
	profiles.New,
))

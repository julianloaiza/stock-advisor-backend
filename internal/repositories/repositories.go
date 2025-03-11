package repositories

import (
	"github.com/julianloaiza/stock-advisor/internal/repositories/profiles"
	"github.com/julianloaiza/stock-advisor/internal/repositories/users"
	"go.uber.org/fx"
)

var Module = fx.Module("repositories", fx.Provide(
	users.New,
	profiles.New,
))

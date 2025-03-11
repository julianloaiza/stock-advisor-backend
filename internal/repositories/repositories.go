package repositories

import (
	"github.com/disturb16/go-service-struct/simple/internal/repositories/profiles"
	"github.com/disturb16/go-service-struct/simple/internal/repositories/users"
	"go.uber.org/fx"
)

var Module = fx.Module("repositories", fx.Provide(
	users.New,
	profiles.New,
))

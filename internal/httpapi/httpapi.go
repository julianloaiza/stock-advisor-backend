package httpapi

import (
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/auth"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers/profiles"
	"go.uber.org/fx"
)

var Module = fx.Module("httpapi", fx.Provide(
	auth.New,
	profiles.New,
))

package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/database"
	"github.com/julianloaiza/stock-advisor/internal/httpapi"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/julianloaiza/stock-advisor/internal/repositories"
	"github.com/julianloaiza/stock-advisor/internal/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Lc       fx.Lifecycle
	Config   *config.Config
	DB       *sql.DB
	Echo     *echo.Echo
	Handlers []handlers.Handler `group:"handlers"`
}

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			config.New,
			database.New,
			echo.New,
		),
		repositories.Module,
		services.Module,
		httpapi.Module,
		fx.Invoke(
			setLifeCycle,
		),
	)

	app.Run()
}

func setLifeCycle(p Params) {
	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo)
			}

			go func() {
				p.Echo.Logger.Fatal(p.Echo.Start(p.Config.Address))
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := p.Echo.Shutdown(ctx); err != nil {
				log.Println(err)
			}

			if err := p.DB.Close(); err != nil {
				log.Println(err)
			}

			return nil
		},
	})
}

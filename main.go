package main

import (
	"context"
	"log"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/database"
	"github.com/julianloaiza/stock-advisor/internal/httpapi"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/middleware"
	"github.com/julianloaiza/stock-advisor/internal/repositories"
	"github.com/julianloaiza/stock-advisor/internal/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Params inyecta dependencias en el ciclo de vida de la aplicación.
type Params struct {
	fx.In

	Lc       fx.Lifecycle
	Config   *config.Config
	DB       *gorm.DB
	Echo     *echo.Echo
	Handlers []handlers.Handler `group:"handlers"`
}

// main inicia la aplicación con Uber FX.
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
		fx.Invoke(setLifeCycle),
	)

	app.Run()
}

// setLifeCycle configura el servidor y el cierre de la aplicación.
func setLifeCycle(p Params) {
	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Aplicar CORS con configuración del middleware
			middleware.ApplyCORS(p.Echo, p.Config)

			// Registrar rutas de los handlers
			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo)
			}

			// Iniciar el servidor en una gorutina
			go func() {
				if err := p.Echo.Start(p.Config.Address); err != nil {
					p.Echo.Logger.Error("❌ Error iniciando el servidor:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Cierre del servidor
			if err := p.Echo.Shutdown(ctx); err != nil {
				log.Println("Error al detener el servidor:", err)
			}

			// Cerrar conexión a la base de datos
			sqlDB, err := p.DB.DB()
			if err != nil {
				log.Println("Error al obtener la conexión sql.DB:", err)
			} else {
				if err := sqlDB.Close(); err != nil {
					log.Println("Error al cerrar la conexión a la base de datos:", err)
				}
			}
			return nil
		},
	})
}

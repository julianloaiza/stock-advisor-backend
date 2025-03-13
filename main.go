package main

import (
	"context"
	"log"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/database"
	"github.com/julianloaiza/stock-advisor/internal/httpapi"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/julianloaiza/stock-advisor/internal/repositories"
	"github.com/julianloaiza/stock-advisor/internal/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Params es la estructura que se usa para inyectar las dependencias.
type Params struct {
	fx.In

	Lc       fx.Lifecycle
	Config   *config.Config
	DB       *gorm.DB
	Echo     *echo.Echo
	Handlers []handlers.Handler `group:"handlers"`
}

// main es la función principal que inicia la aplicación.
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

// setLifeCycle configura el ciclo de vida de la aplicación.
func setLifeCycle(p Params) {
	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Registro de rutas
			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo)
			}

			// Inicia el servidor en una gorutina separada
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

			// Obtiene el objeto sql.DB y lo cierra (solo si se usó una conexión real o simulada)
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

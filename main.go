// @title Stock Advisor API
// @version 1.0
// @description API para gestionar y consultar datos de acciones bursátiles
// @host localhost:8080
// @BasePath /
// @tag.name Stocks
// @tag.description Operaciones con acciones bursátiles
package main

import (
	"context"
	"log"
	"time"

	"github.com/julianloaiza/stock-advisor/config"
	"github.com/julianloaiza/stock-advisor/database"
	_ "github.com/julianloaiza/stock-advisor/docs"
	"github.com/julianloaiza/stock-advisor/internal/httpapi"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/handlers"
	"github.com/julianloaiza/stock-advisor/internal/httpapi/middleware"
	"github.com/julianloaiza/stock-advisor/internal/repositories"
	"github.com/julianloaiza/stock-advisor/internal/services"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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

			// Agregar ruta para Swagger
			p.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

			// Registrar rutas de los handlers
			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo)
			}

			// Iniciar el servidor en una gorutina
			go func() {
				log.Printf("🚀 Iniciando servidor en %s", p.Config.Address)
				if err := p.Echo.Start(p.Config.Address); err != nil {
					log.Printf("❌ Error iniciando el servidor: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Añadimos un timeout para el shutdown
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Cierre del servidor
			if err := shutdownServer(shutdownCtx, p.Echo); err != nil {
				log.Printf("Error al detener el servidor: %v", err)
			}

			// Cerrar conexión a la base de datos
			if err := closeDatabase(p.DB); err != nil {
				log.Printf("Error al cerrar la base de datos: %v", err)
			}

			log.Println("✅ Aplicación detenida correctamente")
			return nil
		},
	})
}

// shutdownServer detiene el servidor HTTP.
func shutdownServer(ctx context.Context, e *echo.Echo) error {
	log.Println("🛑 Deteniendo servidor HTTP...")
	return e.Shutdown(ctx)
}

// closeDatabase cierra la conexión a la base de datos.
func closeDatabase(db *gorm.DB) error {
	log.Println("🛑 Cerrando conexión a la base de datos...")
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

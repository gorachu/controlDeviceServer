package router

import (
	"controlDeviceServer/internal/config"
	"controlDeviceServer/internal/router/handlers"
	"controlDeviceServer/internal/storage/sqlite"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(storage *sqlite.Storage, log *slog.Logger, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		reqLogger := log.With(
			slog.String("path", c.Request.URL.Path),
			slog.String("method", c.Request.Method),
			slog.String("client_ip", c.ClientIP()),
		)
		c.Set("logger", reqLogger)
		c.Set("db", storage)
		c.Set("cfg", cfg)
		start := time.Now()
		c.Next()
		reqLogger.Info("request completed",
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.Int("errors", len(c.Errors)),
		)
	})
	tables := r.Group("/tables")
	{
		tables.GET("", handlers.NewHandlers(storage, log).Tables)
		tables.POST("/current_analyzers", handlers.NewHandlers(storage, log).InsertCurrentAnalyzer)
		tables.POST("/controllers", handlers.NewHandlers(storage, log).InsertController)
		tables.POST("/inputmodules", handlers.NewHandlers(storage, log).InsertInputModule)
		tables.POST("/lcds", handlers.NewHandlers(storage, log).InsertLCD)
		tables.POST("/shields", handlers.NewHandlers(storage, log).InsertShield)
	}

	log.Info("starting HTTP server",
		slog.String("port", cfg.HTTPServer.Address))
	return r
}

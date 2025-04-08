package router

import (
	"controlDeviceServer/internal/config"
	"controlDeviceServer/internal/storage/sqlite"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
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
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers", "Access-Control-Allow-Credentials"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}
	r.Use(cors.New(corsConfig))
	/*
		r.OPTIONS("/*path", func(c *gin.Context) {
			if c.Request.Method == "OPTIONS" {
				c.Header("Access-Control-Allow-Origin", "https://localhost:3000")
				c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With,Access-Control-Allow-Origin,Access-Control-Allow-Methods,Access-Control-Allow-Headers,Access-Control-Allow-Credentials")
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Status(204)
				return
			}
			c.Next()
		})*/

	log.Info("starting HTTP server",
		slog.String("port", cfg.HTTPServer.Address))
	if err := r.Run(cfg.HTTPServer.Address); err != nil {
		log.Error("failed to start HTTP server", "error", err)
		return nil
	}
	return r
}

// Package main is the entry point for the backend server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"seismic-monitor/backend/internal/auth"
	"seismic-monitor/backend/internal/config"
	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/ingest"
	"seismic-monitor/backend/internal/api/handlers"
	"seismic-monitor/backend/internal/api/middleware"
	"seismic-monitor/backend/internal/adapters/usgs"

	"github.com/gin-gonic/gin"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	ctx := context.Background()

	cfg := config.Load()

	// 1. Inicializar la conexión a la base de datos
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("No se pudo conectar a la base de datos", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Conexión a PostgreSQL establecida con éxito")

	// 2. Inicializar repositorios y servicios
	userRepo := database.NewUserRepository(db)
	earthquakeRepo := database.NewEarthquakeRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// 3. Inicializar handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userRepo)
	earthquakeHandler := handlers.NewEarthquakeHandler(earthquakeRepo)

	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Habilitar CORS si es necesario (puedes añadir middleware aquí más tarde)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiV1 := r.Group("/api/v1")
	{
		// Rutas públicas
		apiV1.GET("/earthquakes", earthquakeHandler.GetEarthquakes)

		users := apiV1.Group("/users")
		{
			users.POST("/register", authHandler.Register)
			users.POST("/login", authHandler.Login)
		}

		// Rutas protegidas
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.PUT("/users/location", userHandler.UpdateLocation)
		}
	}

	api := r.Group("/api")
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from the API"})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stopWorker := make(chan bool)

	usgsURL := "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.geojson"
	provider := &usgs.USGSAdapter{URL: usgsURL}

	go ingest.StartIngestionWorker(60*time.Second, stopWorker, provider)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	stopWorker <- true

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}

	logger.Info("server stopped")
}

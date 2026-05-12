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

	"seismic-monitor/backend/internal/adapters/email"
	"seismic-monitor/backend/internal/adapters/usgs"
	"seismic-monitor/backend/internal/api/handlers"
	"seismic-monitor/backend/internal/api/middleware"
	"seismic-monitor/backend/internal/auth"
	"seismic-monitor/backend/internal/config"
	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/ingest"
	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/ports"
	"seismic-monitor/backend/internal/services"

	"github.com/gin-gonic/gin"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	ctx := context.Background()

	cfg := config.Load()

	logger.Info("Intentando conectar a la BD...", "url", cfg.DatabaseURL)

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
	earthquakeService := services.NewEarthquakeService(earthquakeRepo)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	reportRepo := &database.ReportRepository{DB: db}

	// 3. Inicializar handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userRepo)
	earthquakeHandler := handlers.NewEarthquakeHandler(earthquakeService)
	reportHandler := &handlers.ReportHandler{Repo: reportRepo}

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
		apiV1.GET("/earthquakes/history", earthquakeHandler.GetHistory)

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
	{
		api.POST("/report-feeling", reportHandler.HandleReport)
	}
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

	usgsURL := "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson"
	provider := &usgs.USGSAdapter{URL: usgsURL}
	var spatialProvider ports.SpatialRepository = userRepo

	// Emails

	// 1. Creamos la "Cola" en memoria (buffer de 100 mensajes)
	alertQueue := make(chan models.AlertMessage, 100)

	// 2. Instanciamos nuestro Adaptador de Emails
	emailAdapter := &email.SMTPSender{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
	}

	// 3. Arrancamos el Worker de Notificaciones (Consumidor)
	go services.StartNotificationWorker(alertQueue, emailAdapter)

	// 4. Instanciamos el Worker de Ingesta
	ingestionWorker := ingest.NewIngestionWorker(
		60*time.Second,
		provider,
		spatialProvider,
		earthquakeRepo,
		alertQueue,
	)

	// Arrancamos el worker en una goroutine
	go ingestionWorker.Start(stopWorker)

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

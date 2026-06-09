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

	"seismic-monitor/backend/internal/adapters/ai"
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

	logger.Info("Attempting to connect to the DB...", "url", cfg.DatabaseURL)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Could not connect to the database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("PostgreSQL connection established successfully")

	alertQueue := make(chan models.AlertMessage, 100)

	userRepo := database.NewUserRepository(db)
	earthquakeRepo := database.NewEarthquakeRepository(db)
	earthquakeService := services.NewEarthquakeService(earthquakeRepo)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	reportRepo := &database.ReportRepository{DB: db}

	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userRepo)
	earthquakeHandler := handlers.NewEarthquakeHandler(earthquakeService)

	geminiKey := os.Getenv("GEMINI_API_KEY")
	aiProvider := &ai.GeminiAdapter{APIKey: geminiKey}

	reportHandler := handlers.NewReportHandler(reportRepo, userRepo, alertQueue)

	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/earthquakes", earthquakeHandler.GetEarthquakes)
		apiV1.GET("/earthquakes/history", earthquakeHandler.GetHistory)

		users := apiV1.Group("/users")
		{
			users.POST("/register", authHandler.Register)
			users.POST("/login", authHandler.Login)
		}

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

	emailAdapter := &email.SMTPSender{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
	}

	go services.StartNotificationWorker(ctx, alertQueue, emailAdapter, aiProvider)

	services.StartReportCleanupWorker(reportRepo)

	ingestionWorker := ingest.NewIngestionWorker(
		60*time.Second,
		provider,
		spatialProvider,
		earthquakeRepo,
		alertQueue,
	)

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

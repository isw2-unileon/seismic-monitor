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

	"seismic-monitor/internal/database"
	"seismic-monitor/internal/ingest"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/config"
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

	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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

	ingestTask := func() {
		url := "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.geojson"
		logger.Info("Iniciando polling al USGS...")

		data, err := ingest.FetchData(url)
		if err != nil {
			logger.Error("Error al obtener datos", "error", err)
			return
		}

		response, err := ingest.ParseUSGSData(data)
		if err != nil {
			logger.Error("Error al parsear datos", "error", err)
			return
		}

		logger.Info("Datos procesados correctamente", "sismos_detectados", len(response.Features))

		// Imprimimos el primer sismo solo para verificar en los logs que funciona
		if len(response.Features) > 0 {
			primerSismo := response.Features[0]
			logger.Info("Último sismo detectado",
				"id", primerSismo.ID,
				"lugar", primerSismo.Info.Place,
				"magnitud", primerSismo.Info.Mag,
			)
		}

		// TODO: Guardar en PostgreSQL
	}

	go ingest.StartIngestionWorker(60*time.Second, stopWorker, ingestTask)

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

package ingest

import (
	"log"
	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/ports"
	"time"
)

type IngestionWorker struct {
	interval       time.Duration
	provider       ports.EarthquakeProvider
	spatialRepo    ports.SpatialRepository
	earthquakeRepo ports.EarthquakeRepository
	alertQueue     chan<- models.AlertMessage
}

func NewIngestionWorker(
	interval time.Duration,
	provider ports.EarthquakeProvider,
	spatialRepo ports.SpatialRepository,
	earthquakeRepo ports.EarthquakeRepository,
	alertQueue chan<- models.AlertMessage,
) *IngestionWorker {
	return &IngestionWorker{
		interval:       interval,
		provider:       provider,
		spatialRepo:    spatialRepo,
		earthquakeRepo: earthquakeRepo,
		alertQueue:     alertQueue,
	}
}

func (w *IngestionWorker) Start(stopChan <-chan bool) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Motor de ingesta iniciado cada %v\n", w.interval)

	for {
		select {
		case <-ticker.C:
			// Nota que ahora usamos "w.provider", "w.earthquakeRepo", etc.
			response, _ := w.provider.GetEarthquakes()
			for _, sismo := range response.Features {

				err := w.earthquakeRepo.SaveEarthquake(sismo)
				if err != nil {
					log.Printf("[Worker] Aviso al guardar sismo %s: %v", sismo.ID, err)
				}

				affectedUsers, _ := w.spatialRepo.GetAffectedUsers(sismo)
				for _, user := range affectedUsers {
					w.alertQueue <- models.AlertMessage{
						User:  user,
						Sismo: sismo,
					}
				}
			}
		case <-stopChan:
			log.Println("[Worker] Deteniendo el motor de ingesta...")
			return
		}
	}
}

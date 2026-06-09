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
	processedIDs   map[string]bool
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
		processedIDs:   make(map[string]bool),
	}
}

func (w *IngestionWorker) processEarthquakes() {
	response, err := w.provider.GetEarthquakes()
	if err != nil {
		log.Printf("[Worker] Error fetching earthquakes: %v", err)
		return
	}

	newSavedCount := 0

	for _, sismo := range response.Features {

		if w.processedIDs[sismo.ID] {
			continue
		}

		err := w.earthquakeRepo.SaveEarthquake(sismo)
		if err != nil {
			log.Printf("[Worker] Warning saving earthquake %s: %v", sismo.ID, err)
			continue
		}

		w.processedIDs[sismo.ID] = true
		newSavedCount++

		affectedUsers, _ := w.spatialRepo.GetAffectedUsers(sismo)
		for _, user := range affectedUsers {
			select {
			case w.alertQueue <- models.AlertMessage{
				User:  user,
				Sismo: sismo,
			}:
			default:
				log.Printf("[Worker] Alert queue full, skipping notification for %s", user.Email)
			}
		}
	}

	if newSavedCount > 0 {
		log.Printf("[Worker] SUCCESS: %d new earthquakes saved to database.", newSavedCount)
	}
}

func (w *IngestionWorker) Start(stopChan <-chan bool) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Ingestion motor started every %v\n", w.interval)

	log.Println("[Worker] Running initial earthquake load...")
	w.processEarthquakes()

	for {
		select {
		case <-ticker.C:
			w.processEarthquakes()
		case <-stopChan:
			log.Println("[Worker] Stopping ingestion motor...")
			return
		}
	}
}

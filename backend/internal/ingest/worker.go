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
	processedIDs   map[string]bool // NEW: Map to store already processed earthquake IDs
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
		processedIDs:   make(map[string]bool), // NEW: Initialize the map
	}
}

func (w *IngestionWorker) processEarthquakes() {
	response, err := w.provider.GetEarthquakes()
	if err != nil {
		log.Printf("[Worker] Error fetching earthquakes: %v", err)
		return
	}

	newlySaved := 0 // Counter to know how many we insert in this pass

	for _, sismo := range response.Features {

		// Have we already seen this earthquake in this session? If so, ignore it.
		if w.processedIDs[sismo.ID] {
			continue
		}

		err := w.earthquakeRepo.SaveEarthquake(sismo)
		if err != nil {
			log.Printf("[Worker] Aviso al guardar sismo %s: %v", sismo.ID, err)
			continue // If saving failed, better not to send the alert yet
		}

		w.processedIDs[sismo.ID] = true
		newlySaved++

		affectedUsers, _ := w.spatialRepo.GetAffectedUsers(sismo)
		for _, user := range affectedUsers {
			select {
			case w.alertQueue <- models.AlertMessage{
				User:  user,
				Sismo: sismo,
			}:
			default:
				log.Printf("[Worker] Cola de alertas llena, omitiendo notificación para %s", user.Email)
			}
		}
	}

	if newlySaved > 0 {
		log.Printf("[Worker] ÉXITO: %d sismos nuevos guardados en la base de datos.", newlySaved)
	}
}

func (w *IngestionWorker) Start(stopChan <-chan bool) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Motor de ingesta iniciado cada %v\n", w.interval)

	// THE KEY! Execute the initial load right after starting
	log.Println("[Worker] Ejecutando carga inicial de sismos de la última hora...")
	w.processEarthquakes()

	for {
		select {
		case <-ticker.C:
			w.processEarthquakes()
		case <-stopChan:
			log.Println("[Worker] Deteniendo el motor de ingesta...")
			return
		}
	}
}

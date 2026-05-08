package ingest

import (
	"log"
	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/ports"
	"time"
)

func StartIngestionWorker(
	interval time.Duration,
	stopChan <-chan bool,
	provider ports.EarthquakeProvider,
	spatialRepo ports.SpatialRepository,
	alertQueue chan<- models.AlertMessage, // NUEVO: Canal de salida para las alertas
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Motor de ingesta iniciado cada %v\n", interval)

	for {
		select {
		case <-ticker.C:
			response, _ := provider.GetEarthquakes()
			for _, sismo := range response.Features {
				affectedUsers, _ := spatialRepo.GetAffectedUsers(sismo)

				// En lugar de hacer print, mandamos el trabajo a la cola (canal)
				for _, user := range affectedUsers {
					alertQueue <- models.AlertMessage{
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

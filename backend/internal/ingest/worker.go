package ingest

import (
	"log"
	"time"
	"seismic-monitor/backend/internal/ports"
	"seismic-monitor/backend/internal/models"
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

log.Println("[Worker] Solicitando datos al proveedor...")
			response, err := provider.GetEarthquakes()
			if err != nil {
				log.Printf("[Worker] Error obteniendo sismos: %v\n", err)
				continue
			}

			// Iteramos sobre los sismos detectados
			for _, sismo := range response.Features {
				// TODO: Aquí irá la llamada para guardar el sismo en la Base de Datos.
				// log.Printf("Guardando sismo %s en DB...", sismo.ID)

				// 1. Buscamos colisiones (Usuarios afectados)
				affectedUsers, err := spatialRepo.GetAffectedUsers(sismo)
				if err != nil {
					log.Printf("[Worker] Error buscando colisiones para el sismo %s: %v\n", sismo.ID, err)
					continue
				}

				// 2. Imprimimos el formato requerido para el Sprint 2
				if len(affectedUsers) > 0 {
					log.Printf("¡Peligro! El sismo %s afecta a %d usuarios", sismo.ID, len(affectedUsers))
					for _, user := range affectedUsers {
						// Scaffolding para Sprint 3
						log.Printf("Usuario %s en peligro por sismo %s", user.Email, sismo.ID)
					}
				}
			}
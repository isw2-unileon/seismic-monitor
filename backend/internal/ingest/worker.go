package ingest

import (
	"log"
	"time"

	"github.com/isw2-unileon/seismic-monitor/backend/internal/ports"
)

// StartIngestionWorker ahora recibe el proveedor de sismos y el repositorio espacial
func StartIngestionWorker(interval time.Duration, stopChan <-chan bool, provider ports.EarthquakeProvider, spatialRepo ports.SpatialRepository) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Motor de ingesta iniciado cada %v\n", interval)

	for {
		select {
		case <-ticker.C:
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

		case <-stopChan:
			log.Println("[Worker] Deteniendo el motor de ingesta...")
			return
		}
	}
}

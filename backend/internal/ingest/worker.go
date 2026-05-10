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
	processedIDs   map[string]bool // NUEVO: Mapa para almacenar los IDs de sismos ya procesados
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
		processedIDs:   make(map[string]bool), // NUEVO: Inicializamos el mapa
	}
}

// NUEVO: Método privado que contiene la lógica que antes estaba en el select
func (w *IngestionWorker) processEarthquakes() {
	response, _ := w.provider.GetEarthquakes()

	nuevosGuardados := 0 // Contador para saber cuántos metemos en esta pasada

	for _, sismo := range response.Features {

		// 1. ¿Ya hemos visto este sismo en esta sesión? Si es así, lo ignoramos.
		if w.processedIDs[sismo.ID] {
			continue
		}

		// 2. Guardamos en Base de Datos
		err := w.earthquakeRepo.SaveEarthquake(sismo)
		if err != nil {
			log.Printf("[Worker] Aviso al guardar sismo %s: %v", sismo.ID, err)
			continue // Si falló al guardar, mejor no mandar la alerta todavía
		}

		// 3. Lo marcamos como procesado para la próxima vez
		w.processedIDs[sismo.ID] = true
		nuevosGuardados++

		// 4. Buscamos usuarios y lanzamos alertas
		affectedUsers, _ := w.spatialRepo.GetAffectedUsers(sismo)
		for _, user := range affectedUsers {
			w.alertQueue <- models.AlertMessage{
				User:  user,
				Sismo: sismo,
			}
		}
	}

	if nuevosGuardados > 0 {
		log.Printf("[Worker] ÉXITO: %d sismos nuevos guardados en la base de datos.", nuevosGuardados)
	}
}

func (w *IngestionWorker) Start(stopChan <-chan bool) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Motor de ingesta iniciado cada %v\n", w.interval)

	// ¡LA CLAVE! Ejecutamos la carga inicial nada más arrancar
	log.Println("[Worker] Ejecutando carga inicial de sismos de la última hora...")
	w.processEarthquakes()

	for {
		select {
		case <-ticker.C:
			// En cada tick, volvemos a llamar a la función
			w.processEarthquakes()
		case <-stopChan:
			log.Println("[Worker] Deteniendo el motor de ingesta...")
			return
		}
	}
}

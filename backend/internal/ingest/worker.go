package ingest

import (
	"fmt"
	"seismic-monitor/backend/internal/ports"
	"time"
)

// StartIngestionWorker inicia un bucle infinito en segundo plano
func StartIngestionWorker(interval time.Duration, stopChan <-chan bool, provider ports.EarthquakeProvider, repo ports.EarthquakeRepository) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("Motor de ingesta iniciado cada %v\n", interval)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Buscando nuevos sismos...")
			response, err := provider.GetEarthquakes()
			if err != nil {
				fmt.Printf("Error obteniendo datos: %v\n", err)
				continue
			}
			
			for _, eq := range response.Features {
				err := repo.SaveEarthquake(eq)
				if err != nil {
					fmt.Printf("Error guardando sismo %s: %v\n", eq.ID, err)
				}
			}
			fmt.Printf("Se procesaron %d sismos exitosamente.\n", len(response.Features))

		case <-stopChan:
			fmt.Println("Deteniendo el motor de ingesta...")
			return
		}
	}
}

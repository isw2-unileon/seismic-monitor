package ingest

import (
	"fmt"
	"seismic-monitor/internal/ports"
	"time"
)

// StartIngestionWorker inicia un bucle infinito en segundo plano
func StartIngestionWorker(interval time.Duration, stopChan <-chan bool, provider ports.EarthquakeProvider) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("Motor de ingesta iniciado cada %v\n", interval)

	for {
		select {
		case <-ticker.C:
			// Llamamos al proveedor sin saber si es el USGS, otra API o un mock de testing
			fmt.Println("Buscando nuevos sismos...")
			response, err := provider.GetEarthquakes()
			if err != nil {
				fmt.Printf("Error obteniendo datos: %v\n", err)
				continue
			}
			fmt.Printf("Se procesaron %d sismos exitosamente.\n", len(response.Features))
			// TODO: Guardar response.Features en Base de Datos

		case <-stopChan:
			fmt.Println("Deteniendo el motor de ingesta...")
			return
		}
	}
}

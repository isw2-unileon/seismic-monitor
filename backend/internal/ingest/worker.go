package ingest

import (
	"fmt"
	"time"
)

// StartIngestionWorker inicia un bucle infinito en segundo plano
func StartIngestionWorker(interval time.Duration, stopChan <-chan bool, task func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("Motor de ingesta iniciado cada %v\n", interval)

	for {
		select {
		case <-ticker.C:
			task()
		case <-stopChan:
			fmt.Println("Deteniendo el motor de ingesta...")
			return
		}
	}
}

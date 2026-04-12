package ingest

import (
	"testing"
	"time"
)

func TestStartIngestionWorker(t *testing.T) {
	stopChan := make(chan bool)

	executed := false
	mockTask := func() {
		executed = true
		stopChan <- true
	}

	go StartIngestionWorker(100*time.Millisecond, stopChan, mockTask)

	select {
	case <-time.After(500 * time.Millisecond):
		t.Error("El worker no se ejecutó a tiempo")
	case <-stopChan:
		if !executed {
			t.Error("La tarea no fue llamada por el worker")
		}
	}
}

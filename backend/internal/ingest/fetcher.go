package ingest

import (
	"fmt"
	"io"
	"net/http"
)

// FetchData realiza una petición GET a la URL de USGS y devuelve los bytes
func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con USGS: %w", err)
	}
	// Muy importante en Go: liberar el recurso al terminar la función
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("código de estado no esperado: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer el cuerpo de la respuesta: %w", err)
	}

	return body, nil
}

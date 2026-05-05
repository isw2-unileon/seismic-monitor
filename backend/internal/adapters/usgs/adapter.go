package usgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"seismic-monitor/backend/internal/models"
)

// USGSAdapter es el adaptador que se conecta al exterior.
type USGSAdapter struct {
	URL string
}

// GetEarthquakes cumple con la interfaz ports.EarthquakeProvider
func (a *USGSAdapter) GetEarthquakes() (models.USGSResponse, error) {
<<<<<<< HEAD
	data, err := FetchData(a.URL)
=======
	data, err := a.fetchData(a.URL)
>>>>>>> b1ac3b915bac750d7595dcce04e4d70208896d44
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en fetch: %w", err)
	}

<<<<<<< HEAD
	response, err := ParseUSGSData(data)
=======
	response, err := a.parseUSGSData(data)
>>>>>>> b1ac3b915bac750d7595dcce04e4d70208896d44
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en parseo: %w", err)
	}

	return response, nil
}

<<<<<<< HEAD
// FetchData realiza la petición HTTP GET a la URL proporcionada
func FetchData(url string) ([]byte, error) {
=======
func (a *USGSAdapter) fetchData(url string) ([]byte, error) {
>>>>>>> b1ac3b915bac750d7595dcce04e4d70208896d44
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
<<<<<<< HEAD
		return nil, fmt.Errorf("bad status: %s", resp.Status)
=======
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
>>>>>>> b1ac3b915bac750d7595dcce04e4d70208896d44
	}

	return io.ReadAll(resp.Body)
}

<<<<<<< HEAD
// ParseUSGSData deserializa el JSON en la estructura USGSResponse
func ParseUSGSData(data []byte) (models.USGSResponse, error) {
	var response models.USGSResponse
	err := json.Unmarshal(data, &response)
	return response, err
=======
func (a *USGSAdapter) parseUSGSData(data []byte) (models.USGSResponse, error) {
	var response models.USGSResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return models.USGSResponse{}, err
	}
	return response, nil
>>>>>>> b1ac3b915bac750d7595dcce04e4d70208896d44
}

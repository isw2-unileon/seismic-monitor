package models

// USGSResponse es la estructura raíz para el GeoJSON completo
type USGSResponse struct {
	Features []Earthquake `json:"features"`
}

// Earthquake representa un evento sísmico individual
type Earthquake struct {
	ID       string             `json:"id"`
	Info     EarthquakeProps    `json:"properties"`
	Geometry EarthquakeGeometry `json:"geometry"`
}

// EarthquakeProps contiene los detalles del sismo
type EarthquakeProps struct {
	Mag   float64 `json:"mag"`
	Place string  `json:"place"`
	Time  int64   `json:"time"` // El tiempo viene en milisegundos Unix
}

// EarthquakeGeometry contiene las coordenadas [longitud, latitud, profundidad]
type EarthquakeGeometry struct {
	Coordinates []float64 `json:"coordinates"`
}

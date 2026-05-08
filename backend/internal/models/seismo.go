package models

// USGSResponse es el objeto raíz que recibimos de la API
type USGSResponse struct {
	Type     string    `json:"type"` // Suele ser "FeatureCollection"
	Features []Feature `json:"features"`
}

// Feature representa un sismo individual.
// Es lo que antes llamábamos "Earthquake"
type Feature struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"` // Suele ser "Feature"
	Info     EarthquakeProps    `json:"properties"`
	Geometry EarthquakeGeometry `json:"geometry"`
}

// EarthquakeProps contiene los detalles descriptivos
type EarthquakeProps struct {
	Mag   float64 `json:"mag"`
	Place string  `json:"place"`
	Time  int64   `json:"time"` // Tiempo en milisegundos Unix
}

// EarthquakeGeometry contiene los datos espaciales
type EarthquakeGeometry struct {
	Type        string    `json:"type"`        // "Point"
	Coordinates []float64 `json:"coordinates"` // [longitud, latitud, profundidad]
}

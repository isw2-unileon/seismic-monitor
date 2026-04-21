package models

// USGSResponse es la estructura raíz para el GeoJSON completo
type USGSResponse struct {
	Features []Earthquake `json:"features"`
}

// Earthquake representa un evento sísmico individual (Feature en GeoJSON)
type Earthquake struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"` // "Feature"
	Info     EarthquakeProps    `json:"properties"`
	Geometry EarthquakeGeometry `json:"geometry"`
}

// FeatureCollection es el contenedor raíz para GeoJSON
type FeatureCollection struct {
	Type     string       `json:"type"` // "FeatureCollection"
	Features []Earthquake `json:"features"`
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

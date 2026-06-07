package models

// USGSResponse is the root object we receive from the API
type USGSResponse struct {
	Type     string    `json:"type"` // Usually "FeatureCollection"
	Features []Feature `json:"features"`
}

// Feature represents an individual earthquake.
// It is what we used to call "Earthquake"
type Feature struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"` // Usually "Feature"
	Info     EarthquakeProps    `json:"properties"`
	Geometry EarthquakeGeometry `json:"geometry"`
	AIAdvice string             `json:"ai_advice,omitempty"`
}

// EarthquakeProps contains the descriptive details
type EarthquakeProps struct {
	Mag   float64 `json:"mag"`
	Place string  `json:"place"`
	Time  int64   `json:"time"` // Time in Unix milliseconds
}

// EarthquakeGeometry contains the spatial data
type EarthquakeGeometry struct {
	Type        string    `json:"type"`        // "Point"
	Coordinates []float64 `json:"coordinates"` // [longitude, latitude, depth]
}

<script setup>
import { onMounted, onUnmounted, shallowRef } from 'vue'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import { apiService } from '../services/api'

// shallowRef is mandatory here. It prevents Vue's reactivity system
// from proxying Leaflet's massive internal object, avoiding severe memory leaks.
const mapContainer = shallowRef(null)
const mapInstance = shallowRef(null)

onMounted(async () => {
  // 1. Defensive check: Ensure the DOM element exists before attaching Leaflet
  if (!mapContainer.value) {
    console.error('Map container is not ready in the DOM.')
    return
  }

  // 2. Initialize the map centered near León, Spain
  mapInstance.value = L.map(mapContainer.value).setView([42.5987, -5.5671], 5)

  // 3. Attach OpenStreetMap base tiles
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap',
  }).addTo(mapInstance.value)

  // 4. Asynchronously fetch GeoJSON data from our mock API service
  try {
    const geoJsonData = await apiService.getEarthquakes()

    // 5. Parse and render the GeoJSON features on the map
    L.geoJSON(geoJsonData, {
      // Transform default markers into CircleMarkers to reflect earthquake magnitude
      pointToLayer: (feature, latlng) => {
        const magnitude = feature.properties.magnitude || 1
        return L.circleMarker(latlng, {
          radius: magnitude * 3, // Visual scaling based on magnitude
          fillColor: '#ff4444',
          color: '#222',
          weight: 1,
          opacity: 1,
          fillOpacity: 0.7,
        })
      },
      // Attach interactive popups to each point
      onEachFeature: (feature, layer) => {
        if (feature.properties) {
          const time = new Date(feature.properties.time).toLocaleString()
          layer.bindPopup(`
            <div style="font-family: sans-serif;">
              <strong style="color: red;">Magnitude: ${feature.properties.magnitude}</strong><br>
              <strong>Location:</strong> ${feature.properties.place}<br>
              <strong>Depth:</strong> ${feature.properties.depth_km} km<br>
              <strong>Time:</strong> ${time}
            </div>
          `)
        }
      },
    }).addTo(mapInstance.value)
  } catch (error) {
    // In a real scenario, this should trigger a UI toast/alert, not just a console log
    console.error('Critical failure fetching seismic data:', error)
  }
})

onUnmounted(() => {
  // Strict memory management: destroy Leaflet instance when navigating away from the view
  if (mapInstance.value) {
    mapInstance.value.remove()
  }
})
</script>

<template>
  <div class="map-wrapper">
    <div ref="mapContainer" class="seismic-map"></div>
  </div>
</template>

<style scoped>
.map-wrapper {
  width: 100%;
  height: 100vh; /* Takes full viewport height */
  display: flex;
  flex-direction: column;
}

.seismic-map {
  flex-grow: 1;
  width: 100%;
  z-index: 1; /* Ensures map stays below fixed headers/modals if added later */
}
</style>

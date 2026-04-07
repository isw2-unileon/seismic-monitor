<script setup>
import { onMounted, onUnmounted, shallowRef } from 'vue'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import { apiService } from '../services/api'

// shallowRef prevents Vue's reactivity system from proxying Leaflet's internal object, avoiding memory leaks.
const mapContainer = shallowRef(null)
const mapInstance = shallowRef(null)

onMounted(async () => {
  if (!mapContainer.value) {
    console.error('Map container is not ready in the DOM.')
    return
  }

  const earthBounds = L.latLngBounds([
    [-90, -180],
    [90, 180],
  ])

  mapInstance.value = L.map(mapContainer.value, {
    center: [42.5987, -5.5671],
    zoom: 5,
    minZoom: 3,
    maxBounds: earthBounds,
    maxBoundsViscosity: 1.0,
  })

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    noWrap: true, // Prevents infinite horizontal panning
    bounds: earthBounds,
    attribution: '© OpenStreetMap',
  }).addTo(mapInstance.value)

  try {
    const geoJsonData = await apiService.getEarthquakes()

    L.geoJSON(geoJsonData, {
      pointToLayer: (feature, latlng) => {
        const magnitude = feature.properties.magnitude || 1
        return L.circleMarker(latlng, {
          radius: magnitude * 3,
          fillColor: '#ff4444',
          color: '#222',
          weight: 1,
          opacity: 1,
          fillOpacity: 0.7,
        })
      },
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
    // TODO: Replace with UI toast/alert in production
    console.error('Critical failure fetching seismic data:', error)
  }
})

onUnmounted(() => {
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
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.seismic-map {
  flex-grow: 1;
  width: 100%;
  z-index: 1; /* Keeps map below future fixed headers/modals */
}
</style>

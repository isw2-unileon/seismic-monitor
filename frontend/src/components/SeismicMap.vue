<script setup>
import { ref, onMounted, onUnmounted, shallowRef } from 'vue'
import { useRouter } from 'vue-router'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import { apiService } from '../services/api'

const router = useRouter()
const isMenuOpen = ref(false)
const mapContainer = shallowRef(null)
const mapInstance = shallowRef(null)

const layersGroup = L.layerGroup()
const tempLayer = L.layerGroup()

let tempCircle = null
const tempRadius = ref(100)
const tempMagnitude = ref(3.0)
const pendingLocation = ref(null)

const selectedMarkerId = ref(null) 
const showDeleteConfirm = ref(false) 
let resizeListener = null // Referencia para el evento de redimensión

const getIconUrl = (name) => new URL(`../assets/icons/${name}.png`, import.meta.url).href
const toggleMenu = () => isMenuOpen.value = !isMenuOpen.value
const navigateTo = (routeName) => { isMenuOpen.value = false; router.push({ name: routeName }); }
const handleLogout = () => { localStorage.removeItem('auth_token'); localStorage.removeItem('user_data'); router.push({ name: 'login' }); }

const customMarkerIcon = L.icon({
  iconUrl: getIconUrl('marker'), // Apunta a src/assets/icons/marker.png
  iconSize: [48, 48],            // Tamaño en el que se pintará [ancho, alto]
  iconAnchor: [24, 48],          // Punto del icono que ancla a la lat/lng (centro-base)
  popupAnchor: [0, -64]          // Punto desde el que emerge el popup (centro-arriba)
})


const loadUserCenters = () => {
  const data = JSON.parse(localStorage.getItem('user_data') || '{}')
  let centers = data.alert_centers || []
  
  // Si no hay centros pero hay ubicación base, la reconstruimos como el primer centro
  if (centers.length === 0 && data.latitude && data.longitude) {
    centers = [{
      id: 'base',
      lat: data.latitude,
      lng: data.longitude,
      radius: data.alert_radius || data.alert_radius_km || 100,
      min_magnitude: data.min_magnitude || 3.0
    }]
  }

  return {
    centers: centers,
    defaultRadius: data.alert_radius || data.alert_radius_km || 100,
    defaultMagnitude: data.min_magnitude || 3.0
  }
}

const renderAllCenters = () => {
  layersGroup.clearLayers()
  const { centers } = loadUserCenters()

  centers.forEach(center => {
    const marker = L.marker([center.lat, center.lng], { icon: customMarkerIcon }).addTo(layersGroup)
    
    marker.on('click', (e) => {
      L.DomEvent.stopPropagation(e)
      cancelLocation()
      selectedMarkerId.value = center.id
      showDeleteConfirm.value = false
    })

    L.circle([center.lat, center.lng], {
      radius: center.radius * 1000,
      color: '#e94560',
      fillColor: '#e94560',
      fillOpacity: 0.15,
      weight: 2,
      dashArray: '5, 10'
    }).addTo(layersGroup)
  })
}

const confirmLocation = async () => {
  if (!pendingLocation.value) return
  const { centers } = loadUserCenters()
  const newCenter = {
    id: Date.now(),
    lat: pendingLocation.value.lat,
    lng: pendingLocation.value.lng,
    radius: tempRadius.value,
    min_magnitude: tempMagnitude.value
  }
  const data = JSON.parse(localStorage.getItem('user_data') || '{}')
  data.alert_centers = [...centers, newCenter]
  data.alert_radius_km = tempRadius.value 
  data.min_magnitude = tempMagnitude.value
  localStorage.setItem('user_data', JSON.stringify(data))
  
  try {
    await apiService.updateUserSettings({
      latitude: newCenter.lat,
      longitude: newCenter.lng,
      alert_radius: newCenter.radius,
      min_magnitude: newCenter.min_magnitude
    })
  } catch (error) {
    console.error("Error sincronizando con el servidor:", error)
  }

  cancelLocation() 
  renderAllCenters()
}

const cancelLocation = () => {
  tempLayer.clearLayers()
  pendingLocation.value = null
  selectedMarkerId.value = null
  if(mapInstance.value) mapInstance.value.closePopup()
}

onMounted(() => {
  if (!mapContainer.value) return
  const { centers, defaultRadius, defaultMagnitude } = loadUserCenters()
  tempRadius.value = defaultRadius
  tempMagnitude.value = defaultMagnitude

  const earthBounds = L.latLngBounds([[-90, -180], [90, 180]])

  mapInstance.value = L.map(mapContainer.value, {
    center: centers.length > 0 ? [centers[0].lat, centers[0].lng] : [42.5987, -5.5671],
    zoom: 5,
    maxBounds: earthBounds, 
    maxBoundsViscosity: 1.0, 
    preferCanvas: true, 
    zoomAnimation: true, // Mantienes tus animaciones
    zoomSnap: 0.1 // CRÍTICO: Permite fracciones de zoom para que el mapa encaje al milímetro
  })

  // CÁLCULO DINÁMICO DEL ZOOM MÍNIMO HORIZONTAL
  const setDynamicMinZoom = () => {
    if (!mapInstance.value || !mapContainer.value) return
    const containerWidth = mapContainer.value.offsetWidth
    // El mundo mide 256px en zoom 0. Calculamos el zoom para que cubra exactamente tu pantalla
    const minZ = Math.log2(containerWidth / 256)
    mapInstance.value.setMinZoom(minZ)
  }

  // Ejecutamos la primera vez y escuchamos cambios de tamaño de ventana
  setDynamicMinZoom()
  resizeListener = () => setDynamicMinZoom()
  window.addEventListener('resize', resizeListener)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    noWrap: true, 
    bounds: earthBounds, 
    updateWhenZooming: false,
    keepBuffer: 2
  }).addTo(mapInstance.value)

  layersGroup.addTo(mapInstance.value)
  tempLayer.addTo(mapInstance.value)

  const earthquakesLayer = L.layerGroup().addTo(mapInstance.value)
  apiService.getEarthquakesHistory().then(data => {
    if (data && data.features) {
      data.features.forEach(eq => {
        const coords = [eq.geometry.coordinates[1], eq.geometry.coordinates[0]]
        const mag = eq.properties.magnitude
        
        let color = '#4cd137'
        if (mag >= 5) color = '#e84118'
        else if (mag >= 3) color = '#fbc531'

        const circle = L.circleMarker(coords, {
          radius: Math.max(mag * 3, 5),
          fillColor: color,
          color: '#fff',
          weight: 1,
          opacity: 1,
          fillOpacity: 0.8
        }).addTo(earthquakesLayer)

        const time = new Date(eq.properties.time).toLocaleString()
        circle.bindPopup(`
          <div style="font-family: sans-serif; text-align: center;">
            <strong style="color: ${color}; font-size: 16px;">M ${mag.toFixed(1)}</strong><br>
            <span style="font-size: 12px; color: #555;">${time}</span><br>
            <div style="margin-top: 5px; font-size: 14px;">${eq.properties.place}</div>
          </div>
        `)
      })
    }
  }).catch(console.error)

  mapInstance.value.on('click', (e) => {
    selectedMarkerId.value = null 
    tempLayer.clearLayers()
    pendingLocation.value = e.latlng

    const data = JSON.parse(localStorage.getItem('user_data') || '{}')
    tempRadius.value = data.alert_radius_km || 100
    tempMagnitude.value = data.min_magnitude || 3.0

    tempCircle = L.circle(e.latlng, {
      radius: tempRadius.value * 1000,
      color: '#e94560',
      fillColor: '#e94560',
      fillOpacity: 0.3,
      weight: 2,
      dashArray: '5, 10'
    }).addTo(tempLayer)

    const marker = L.marker(e.latlng, { icon: customMarkerIcon }).addTo(tempLayer)

    const popupContent = `
      <div style="min-width: 190px; text-align: center; font-family: sans-serif;">
        <h4 style="margin: 0 0 10px 0; color: #1a1a2e;">Configurar Área</h4>
        
        <div style="margin-bottom: 15px;">
          <div style="display: flex; justify-content: center; align-items: center; gap: 5px; margin-bottom: 5px;">
            <span style="font-size: 14px; color: #333; font-weight: bold;">Radio:</span>
            <input type="number" id="radius-input" min="1" max="5000" step="1" value="${tempRadius.value}" 
                   style="width: 70px; padding: 4px; border: 1px solid #ccc; border-radius: 4px; text-align: center; font-size: 14px; color: #1a1a2e; outline: none;">
            <span style="font-size: 14px; color: #333;">km</span>
          </div>
          <input type="range" id="radius-slider" min="10" max="1000" step="10" value="${tempRadius.value}" 
                 style="width: 100%; accent-color: #e94560; cursor: pointer;">
        </div>

        <div style="margin-bottom: 15px;">
          <div style="display: flex; justify-content: center; align-items: center; gap: 5px; margin-bottom: 5px;">
            <span style="font-size: 14px; color: #333; font-weight: bold;">Magnitud:</span>
            <input type="number" id="magnitude-input" min="0" max="10" step="0.1" value="${tempMagnitude.value}" 
                   style="width: 70px; padding: 4px; border: 1px solid #ccc; border-radius: 4px; text-align: center; font-size: 14px; color: #1a1a2e; outline: none;">
          </div>
          <input type="range" id="magnitude-slider" min="0" max="10" step="0.1" value="${tempMagnitude.value}" 
                 style="width: 100%; accent-color: #fbc531; cursor: pointer;">
        </div>

        <div style="display: flex; gap: 8px;">
          <button id="save-btn" style="flex: 1; background: #28a745; color: white; border: none; padding: 6px; border-radius: 4px; cursor: pointer; font-weight: bold;">Guardar</button>
          <button id="cancel-btn" style="flex: 1; background: #dc3545; color: white; border: none; padding: 6px; border-radius: 4px; cursor: pointer; font-weight: bold;">✕</button>
        </div>
      </div>
    `
    marker.bindPopup(popupContent, { closeButton: false, closeOnClick: false, autoClose: false }).openPopup()
  })

  mapInstance.value.on('popupopen', () => {
    const radiusSlider = document.getElementById('radius-slider')
    const radiusInput = document.getElementById('radius-input')
    const magnitudeSlider = document.getElementById('magnitude-slider')
    const magnitudeInput = document.getElementById('magnitude-input')
    const saveBtn = document.getElementById('save-btn')
    const cancelBtn = document.getElementById('cancel-btn')

    const syncRadius = (value) => {
      let newRadius = parseInt(value)
      if (isNaN(newRadius) || newRadius < 1) newRadius = 1 
      tempRadius.value = newRadius
      if (tempCircle) tempCircle.setRadius(newRadius * 1000)
      if (radiusSlider && radiusSlider.value !== newRadius.toString()) radiusSlider.value = newRadius
      if (radiusInput && radiusInput.value !== newRadius.toString()) radiusInput.value = newRadius
    }

    const syncMagnitude = (value) => {
      let newMag = parseFloat(value)
      if (isNaN(newMag)) newMag = 3.0
      if (newMag < 0) newMag = 0
      if (newMag > 10) newMag = 10
      tempMagnitude.value = newMag
      if (magnitudeSlider && magnitudeSlider.value !== newMag.toString()) magnitudeSlider.value = newMag
      if (magnitudeInput && magnitudeInput.value !== newMag.toString()) magnitudeInput.value = newMag
    }

    if (radiusSlider) radiusSlider.addEventListener('input', (ev) => syncRadius(ev.target.value))
    if (radiusInput) radiusInput.addEventListener('input', (ev) => syncRadius(ev.target.value))
    
    if (magnitudeSlider) magnitudeSlider.addEventListener('input', (ev) => syncMagnitude(ev.target.value))
    if (magnitudeInput) magnitudeInput.addEventListener('input', (ev) => syncMagnitude(ev.target.value))

    if (saveBtn) saveBtn.addEventListener('click', confirmLocation)
    if (cancelBtn) cancelBtn.addEventListener('click', cancelLocation)
  })

  renderAllCenters()
  setTimeout(() => mapInstance.value.invalidateSize(), 250)
})

const deleteCenter = () => {
  if (!selectedMarkerId.value) return
  const data = JSON.parse(localStorage.getItem('user_data') || '{}')
  data.alert_centers = (data.alert_centers || []).filter(c => c.id !== selectedMarkerId.value)
  localStorage.setItem('user_data', JSON.stringify(data))
  selectedMarkerId.value = null
  showDeleteConfirm.value = false
  renderAllCenters()
}

onUnmounted(() => { 
  if (resizeListener) window.removeEventListener('resize', resizeListener)
  if (mapInstance.value) mapInstance.value.remove() 
})
</script>

<template>
  <div class="map-wrapper">
    <div class="controls-overlay">
      <transition name="slide">
        <div v-if="selectedMarkerId" class="delete-wrapper">
          <div v-if="!showDeleteConfirm" class="action-buttons">
            <button @click="showDeleteConfirm = true" class="delete-trigger-btn">🗑 Eliminar Área</button>
            <button @click="selectedMarkerId = null" class="cancel-btn">✕</button>
          </div>
          <div v-else class="confirm-dropdown">
            <p>¿Borrar este área?</p>
            <div class="confirm-actions">
              <button @click="deleteCenter" class="btn-confirm-yes">Sí, borrar</button>
              <button @click="showDeleteConfirm = false" class="btn-confirm-no">No</button>
            </div>
          </div>
        </div>
      </transition>

      <div class="menu-container">
        <button @click="toggleMenu" class="hamburger-btn">
          <span class="bar"></span><span class="bar"></span><span class="bar"></span>
        </button>

        <transition name="fade">
          <nav v-if="isMenuOpen" class="dropdown-menu">
            <button @click="navigateTo('account')" class="menu-item"><img :src="getIconUrl('account')" class="menu-icon" /> Cuenta</button>
            <button @click="navigateTo('history')" class="menu-item"><img :src="getIconUrl('earthquakes')" class="menu-icon" /> Seísmos</button>
            <div class="menu-divider"></div>
            <button @click="handleLogout" class="menu-item logout"><img :src="getIconUrl('logout')" class="menu-icon" /> Salir</button>
          </nav>
        </transition>
      </div>
    </div>
    <div ref="mapContainer" class="seismic-map"></div>
  </div>
</template>

<style scoped src="../assets/styles/components/SeismicMap.css"></style>
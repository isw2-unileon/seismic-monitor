<script setup>
import { ref, onMounted, onUnmounted, shallowRef } from 'vue'
import { useRouter } from 'vue-router'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'

const router = useRouter()
const isMenuOpen = ref(false)
const mapContainer = shallowRef(null)
const mapInstance = shallowRef(null)

const layersGroup = L.layerGroup()
const tempLayer = L.layerGroup()

// Variables para manejar la creación reactiva fuera de Vue
let tempCircle = null
const tempRadius = ref(100) 
const pendingLocation = ref(null)

const selectedMarkerId = ref(null) 
const showDeleteConfirm = ref(false) 

const getIconUrl = (name) => new URL(`../assets/icons/${name}.png`, import.meta.url).href
const toggleMenu = () => isMenuOpen.value = !isMenuOpen.value
const navigateTo = (routeName) => { isMenuOpen.value = false; router.push({ name: routeName }); }
const handleLogout = () => { localStorage.removeItem('auth_token'); localStorage.removeItem('user_data'); router.push({ name: 'login' }); }

const loadUserCenters = () => {
  const data = JSON.parse(localStorage.getItem('user_data') || '{}')
  return {
    centers: data.alert_centers || [],
    defaultRadius: data.alert_radius_km || 100
  }
}

const renderAllCenters = () => {
  layersGroup.clearLayers()
  const { centers } = loadUserCenters()

  centers.forEach(center => {
    const marker = L.marker([center.lat, center.lng]).addTo(layersGroup)
    
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

const confirmLocation = () => {
  if (!pendingLocation.value) return
  const { centers } = loadUserCenters()
  const newCenter = {
    id: Date.now(),
    lat: pendingLocation.value.lat,
    lng: pendingLocation.value.lng,
    radius: tempRadius.value // Usamos el radio final sincronizado
  }
  const data = JSON.parse(localStorage.getItem('user_data') || '{}')
  data.alert_centers = [...centers, newCenter]
  data.alert_radius_km = tempRadius.value // Actualizamos el radio por defecto para el siguiente
  localStorage.setItem('user_data', JSON.stringify(data))
  
  cancelLocation() // Limpia el popup y la capa temporal
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
  const { centers, defaultRadius } = loadUserCenters()
  tempRadius.value = defaultRadius

  // Configuración del mapa preservando tus preferencias de animación
  mapInstance.value = L.map(mapContainer.value, {
    center: centers.length > 0 ? [centers[0].lat, centers[0].lng] : [42.5987, -5.5671],
    zoom: 5,
    preferCanvas: true, 
    zoomAnimation: true,
  })

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    updateWhenZooming: false,
    keepBuffer: 2
  }).addTo(mapInstance.value)

  layersGroup.addTo(mapInstance.value)
  tempLayer.addTo(mapInstance.value)

  // 1. CLIC EN EL MAPA: Crea el marcador y el popup
  mapInstance.value.on('click', (e) => {
    selectedMarkerId.value = null 
    tempLayer.clearLayers()
    pendingLocation.value = e.latlng

    // Recupera el último radio utilizado
    const data = JSON.parse(localStorage.getItem('user_data') || '{}')
    tempRadius.value = data.alert_radius_km || 100

    // Dibuja el círculo temporal transparente
    tempCircle = L.circle(e.latlng, {
      radius: tempRadius.value * 1000,
      color: '#e94560',
      fillColor: '#e94560',
      fillOpacity: 0.3,
      weight: 2,
      dashArray: '5, 10'
    }).addTo(tempLayer)

    const marker = L.marker(e.latlng).addTo(tempLayer)

    // HTML del Popup actualizado con el input numérico para escritura directa
    const popupContent = `
      <div style="min-width: 190px; text-align: center; font-family: sans-serif;">
        <h4 style="margin: 0 0 10px 0; color: #1a1a2e;">Configurar Área</h4>
        
        <div style="display: flex; justify-content: center; align-items: center; gap: 5px; margin-bottom: 10px;">
          <span style="font-size: 14px; color: #333; font-weight: bold;">Radio:</span>
          <input type="number" id="radius-input" min="1" max="5000" step="1" value="${tempRadius.value}" 
                 style="width: 70px; padding: 4px; border: 1px solid #ccc; border-radius: 4px; text-align: center; font-size: 14px; color: #1a1a2e; outline: none;">
          <span style="font-size: 14px; color: #333;">km</span>
        </div>

        <input type="range" id="radius-slider" min="10" max="1000" step="10" value="${tempRadius.value}" 
               style="width: 100%; accent-color: #e94560; margin-bottom: 15px; cursor: pointer;">
        
        <div style="display: flex; gap: 8px;">
          <button id="save-btn" style="flex: 1; background: #28a745; color: white; border: none; padding: 6px; border-radius: 4px; cursor: pointer; font-weight: bold;">Guardar</button>
          <button id="cancel-btn" style="flex: 1; background: #dc3545; color: white; border: none; padding: 6px; border-radius: 4px; cursor: pointer; font-weight: bold;">✕</button>
        </div>
      </div>
    `

    // closeOnClick: false previene que si fallas al arrastrar el slider, se cierre el popup
    marker.bindPopup(popupContent, { closeButton: false, closeOnClick: false, autoClose: false }).openPopup()
  })

  // 2. CONEXIÓN DE EVENTOS DEL DOM AL POPUP (Sincronización Bidireccional)
  mapInstance.value.on('popupopen', () => {
    const slider = document.getElementById('radius-slider')
    const numberInput = document.getElementById('radius-input')
    const saveBtn = document.getElementById('save-btn')
    const cancelBtn = document.getElementById('cancel-btn')

    // Función unificada para actualizar el estado visual y lógico desde cualquier input
    const syncRadius = (value) => {
      let newRadius = parseInt(value)
      if (isNaN(newRadius) || newRadius < 1) newRadius = 1 // Validamos para no romper Leaflet con radios negativos o nulos

      tempRadius.value = newRadius
      if (tempCircle) tempCircle.setRadius(newRadius * 1000)

      // Evitamos bucles de actualización infinita comprobando el valor antes de setearlo
      if (slider && slider.value !== newRadius.toString()) slider.value = newRadius
      if (numberInput && numberInput.value !== newRadius.toString()) numberInput.value = newRadius
    }

    if (slider) {
      slider.addEventListener('input', (ev) => syncRadius(ev.target.value))
    }
    
    if (numberInput) {
      numberInput.addEventListener('input', (ev) => syncRadius(ev.target.value))
    }

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

onUnmounted(() => { if (mapInstance.value) mapInstance.value.remove() })
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
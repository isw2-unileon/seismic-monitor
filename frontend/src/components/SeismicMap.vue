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

// Intentamos importar los iconos. Si fallan, usaremos strings vacíos
let accountIcon, earthquakesIcon, logoutIcon;
try {
  // Usamos import.meta.glob o rutas relativas simples para mayor compatibilidad
  accountIcon = new URL('../assets/icons/account.png', import.meta.url).href
  earthquakesIcon = new URL('../assets/icons/earthquakes.png', import.meta.url).href
  logoutIcon = new URL('../assets/icons/logout.png', import.meta.url).href
} catch (e) {
  console.warn("Iconos no encontrados en src/assets/icons/. Se mostrará solo texto.")
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const handleLogout = () => {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('user_data')
  router.push({ name: 'login' })
}

const navigateTo = (routeName) => {
  isMenuOpen.value = false
  router.push({ name: routeName })
}

onMounted(async () => {
  if (!mapContainer.value) return

  const earthBounds = L.latLngBounds([[-90, -180], [90, 180]])

  // Inicialización defensiva
  mapInstance.value = L.map(mapContainer.value, {
    center: [42.5987, -5.5671],
    zoom: 5,
    minZoom: 3,
    maxBounds: earthBounds,
    maxBoundsViscosity: 1.0,
  })

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    noWrap: true,
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
      }
    }).addTo(mapInstance.value)
  } catch (error) {
    console.error('Error cargando datos:', error)
  }

  // Ajuste forzado de tamaño
  setTimeout(() => {
    if (mapInstance.value) mapInstance.value.invalidateSize()
  }, 200)
})

onUnmounted(() => {
  if (mapInstance.value) mapInstance.value.remove()
})
</script>

<template>
  <div class="map-wrapper">
    <div class="menu-container">
      <button @click="toggleMenu" class="hamburger-btn">
        <span class="bar"></span>
        <span class="bar"></span>
        <span class="bar"></span>
      </button>

      <transition name="fade">
        <nav v-if="isMenuOpen" class="dropdown-menu">
          <button @click="navigateTo('account')" class="menu-item">
            <img :src="accountIcon" alt="" class="menu-icon" />
            Account
          </button>
          <button @click="navigateTo('history')" class="menu-item">
            <img :src="earthquakesIcon" alt="" class="menu-icon" />
            Earthquakes
          </button>
          <div class="menu-divider"></div>
          <button @click="handleLogout" class="menu-item logout">
            <img :src="logoutIcon" alt="" class="menu-icon" />
            Logout
          </button>
        </nav>
      </transition>
    </div>
    
    <div ref="mapContainer" class="seismic-map"></div>
  </div>
</template>

<style scoped>
.map-wrapper {
  position: relative;
  width: 100vw;
  height: 100vh;
  background-color: #000; /* Fondo negro para detectar si el mapa no carga */
}

.seismic-map {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.menu-container {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 10000; /* Z-index extremadamente alto */
}

.hamburger-btn {
  background: #16213e;
  border: 2px solid #e94560;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.bar {
  width: 24px;
  height: 3px;
  background-color: white;
  display: block;
}

.dropdown-menu {
  position: absolute;
  top: 60px;
  right: 0;
  background: #16213e;
  border: 1px solid #2a3158;
  border-radius: 8px;
  min-width: 200px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.8);
}

.menu-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  font-size: 1rem;
}

.menu-icon {
  width: 20px;
  height: 20px;
  object-fit: contain;
}

.menu-divider {
  height: 1px;
  background: #2a3158;
}

.logout {
  color: #e94560;
}

.fade-enter-active, .fade-leave-active {
  transition: all 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>
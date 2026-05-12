<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../services/api'

const router = useRouter()
const userData = ref({
  name: '',
  email: '',
  min_magnitude: 3.0,
  alert_radius_km: 100,
  latitude: 0,
  longitude: 0,
  alert_centers: [] 
  // Nota: Mantenemos el radius guardado en la lógica interna del localStorage 
  // para que el mapa recuerde el último que usaste, pero ya no se muestra aquí.
})

onMounted(() => {
  const savedData = localStorage.getItem('user_data')
  if (savedData) {
    const parsed = JSON.parse(savedData)
    userData.value = {
      ...userData.value,
      ...parsed,
      alert_centers: parsed.alert_centers || []
    }
  }
})

const saveSettings = async () => {
  localStorage.setItem('user_data', JSON.stringify(userData.value))
  
  try {
    // Si hay centros, usamos el último para la ubicación del backend
    const lastCenter = userData.value.alert_centers.length > 0 
      ? userData.value.alert_centers[userData.value.alert_centers.length - 1]
      : { lat: userData.value.latitude, lng: userData.value.longitude, radius: userData.value.alert_radius_km }

    await apiService.updateUserSettings({
      latitude: lastCenter.lat,
      longitude: lastCenter.lng,
      alert_radius: lastCenter.radius || userData.value.alert_radius_km,
      min_magnitude: userData.value.min_magnitude
    })
    alert('Preferencias actualizadas y sincronizadas')
  } catch (error) {
    console.error("Error sincronizando settings:", error)
    alert('Preferencias guardadas localmente, pero error al sincronizar con el servidor')
  }
}

const removeCenter = (id) => {
  userData.value.alert_centers = userData.value.alert_centers.filter(c => c.id !== id)
  saveSettings()
}

const goBack = () => router.push({ name: 'map' })
</script>

<template>
  <div class="account-wrapper">
    <div class="settings-card">
      <header class="card-header">
        <button @click="goBack" class="back-btn">← Volver al Mapa</button>
        <h1>Configuración de Cuenta</h1>
      </header>

      <form @submit.prevent="saveSettings" class="settings-form">
        <div class="form-section">
          <h3><span class="icon">👤</span> Información Personal</h3>
          <div class="form-group">
            <label>Nombre Completo</label>
            <input v-model="userData.name" type="text" placeholder="Ej: Juan Pérez">
          </div>
          <div class="form-group">
            <label>Email</label>
            <input v-model="userData.email" type="email" disabled class="disabled-input">
          </div>
        </div>

        <div class="form-section">
          <h3><span class="icon">⚙️</span> Preferencias de Alerta</h3>
          <div class="form-group">
            <label>Magnitud Mínima de Alerta (Global)</label>
            <div class="magnitude-control">
              <input v-model.number="userData.min_magnitude" type="range" min="0" max="10" step="0.1" class="mag-slider">
              <span class="mag-value">M {{ userData.min_magnitude }}</span>
            </div>
            <p class="help-text">Solo se te notificará de terremotos con una magnitud igual o superior a este valor.</p>
          </div>
        </div>

        <div class="form-section">
          <h3><span class="icon">📡</span> Mis Zonas de Alerta</h3>
          <div class="centers-list">
            <div v-if="userData.alert_centers.length === 0" class="no-data">
              No tienes zonas guardadas. Haz clic en el mapa para añadir una.
            </div>
            
            <div v-for="center in userData.alert_centers" :key="center.id" class="center-item">
              <div class="center-info">
                <span class="coords">📍 {{ center.lat.toFixed(4) }}, {{ center.lng.toFixed(4) }}</span>
                <div class="badges">
                  <span class="badge radius">{{ center.radius }} km</span>
                  <span class="badge magnitude" v-if="center.min_magnitude">M {{ center.min_magnitude }}</span>
                </div>
              </div>
              <button type="button" @click="removeCenter(center.id)" class="btn-delete-small" title="Eliminar zona">🗑</button>
            </div>
          </div>
        </div>

        <button type="submit" class="save-btn">Guardar Cambios</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.account-wrapper {
  height: 100vh;
  overflow-y: auto;
  background-color: #1a1a2e;
  padding: 2rem 1rem 10rem 1rem; 
  display: flex;
  justify-content: center;
  align-items: flex-start;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

.settings-card { 
  background: #16213e; 
  width: 100%; 
  max-width: 650px; 
  border-radius: 12px; 
  padding: 2.5rem; 
  border: 1px solid #2a3158; 
  color: #fff; 
  box-shadow: 0 10px 30px rgba(0,0,0,0.5); 
}

.card-header { display: flex; align-items: center; gap: 1.5rem; margin-bottom: 2.5rem; }
.back-btn { background: transparent; border: 1px solid #e94560; color: #e94560; padding: 8px 16px; border-radius: 6px; cursor: pointer; transition: 0.3s; }
.back-btn:hover { background: #e94560; color: white; }
.form-section { margin-bottom: 3rem; }
.form-section h3 { color: #e94560; border-bottom: 1px solid #2a3158; padding-bottom: 0.5rem; margin-bottom: 1.5rem; display: flex; align-items: center; gap: 10px; }
.form-group { margin-bottom: 1.2rem; }
.form-group label { display: block; margin-bottom: 0.5rem; color: #a0aab2; font-size: 0.9rem; }
input[type="text"], input[type="email"] { width: 100%; background: #0f172a; border: 1px solid #2a3158; color: white; padding: 12px; border-radius: 6px; }
.disabled-input { opacity: 0.6; cursor: not-allowed; }
.centers-list { display: flex; flex-direction: column; gap: 10px; }
.center-item { display: flex; justify-content: space-between; align-items: center; background: #1f2937; padding: 12px 16px; border-radius: 8px; border-left: 4px solid #e94560; }
.center-info { display: flex; gap: 15px; align-items: center; }
.coords { font-family: 'Courier New', Courier, monospace; color: #fff; }
.badges { display: flex; gap: 8px; }
.badge { background: #2a3158; padding: 2px 8px; border-radius: 4px; font-size: 0.8rem; color: #e94560; }
.badge.magnitude { color: #fbc531; }
.magnitude-control { display: flex; align-items: center; gap: 15px; background: #0f172a; padding: 15px; border-radius: 8px; border: 1px solid #2a3158; }
.mag-slider { flex: 1; accent-color: #fbc531; cursor: pointer; }
.mag-value { font-weight: bold; color: #fbc531; min-width: 60px; text-align: right; font-size: 1.1rem; }
.help-text { font-size: 0.8rem; color: #6b7280; margin-top: 8px; margin-left: 5px; }
.btn-delete-small { background: transparent; border: none; font-size: 1.2rem; cursor: pointer; filter: grayscale(1); transition: 0.2s; }
.btn-delete-small:hover { filter: grayscale(0); transform: scale(1.1); }
.no-data { text-align: center; padding: 2rem; color: #6b7280; border: 2px dashed #2a3158; border-radius: 8px; }
.save-btn { width: 100%; padding: 15px; background: #e94560; color: white; border: none; border-radius: 8px; font-weight: bold; font-size: 1rem; cursor: pointer; }
.save-btn:hover { background: #d63d56; }
</style>
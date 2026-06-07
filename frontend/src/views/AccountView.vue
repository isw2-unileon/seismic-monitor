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
  // Note: We keep the radius saved in the internal localStorage logic
  // so the map remembers the last one you used, but it is no longer shown here.
})

onMounted(() => {
  const savedData = localStorage.getItem('user_data')
  if (savedData) {
    const parsed = JSON.parse(savedData)
    userData.value = {
      ...userData.value,
      ...parsed,
      // Map backend 'alert_radius' to frontend 'alert_radius_km' if needed
      alert_radius_km: parsed.alert_radius || parsed.alert_radius_km || 100,
      alert_centers: parsed.alert_centers || []
    }
  }
})

const saveSettings = async () => {
  console.log("Saving user settings:", userData.value)
  try {
    // If there are centers, we use the last one for the backend location
    const lastCenter = userData.value.alert_centers.length > 0 
      ? userData.value.alert_centers[userData.value.alert_centers.length - 1]
      : { lat: userData.value.latitude, lng: userData.value.longitude, radius: userData.value.alert_radius_km }

    const payload = {
      name: userData.value.name,
      latitude: lastCenter.lat,
      longitude: lastCenter.lng,
      alert_radius: lastCenter.radius || userData.value.alert_radius_km,
      min_magnitude: userData.value.min_magnitude
    }
    console.log("Sending payload to backend:", payload)

    await apiService.updateUserSettings(payload)

    localStorage.setItem('user_data', JSON.stringify(userData.value))
    alert('Preferences updated and synchronized')
  } catch (error) {
    console.error("Error synchronizing settings:", error)
    alert('Error synchronizing with the server: ' + error.message)
  }
}

const removeCenter = async (id) => {
  try {
    if (typeof id === 'string') {
      await apiService.deleteUserLocation(id)
    }
  } catch (error) {
    console.error("Error deleting center:", error)
  }
  userData.value.alert_centers = userData.value.alert_centers.filter(c => c.id !== id)
  saveSettings()
}


const goBack = () => router.push({ name: 'map' })
</script>

<template>
  <div class="account-wrapper">
    <div class="settings-card">
      <header class="card-header">
        <button @click="goBack" class="back-btn">← Back to Map</button>
        <h1>Account Settings</h1>
      </header>

      <form @submit.prevent="saveSettings" class="settings-form">
        <div class="form-section">
          <h3><span class="icon">👤</span> Personal Information</h3>
          <div class="form-group">
            <label>Full Name</label>
            <input v-model="userData.name" type="text" placeholder="e.g.: John Doe">
          </div>
          <div class="form-group">
            <label>Email</label>
            <input v-model="userData.email" type="email" disabled class="disabled-input">
          </div>
        </div>

        <div class="form-section">
          <h3><span class="icon">⚙️</span> Alert Preferences</h3>
          <div class="form-group">
            <label>Minimum Alert Magnitude (Global)</label>
            <div class="magnitude-control">
              <input v-model.number="userData.min_magnitude" type="range" min="0" max="10" step="0.1" class="mag-slider">
              <span class="mag-value">M {{ userData.min_magnitude }}</span>
            </div>
            <p class="help-text">You will only be notified of earthquakes with a magnitude equal to or greater than this value.</p>
          </div>
        </div>

        <div class="form-section">
          <h3><span class="icon">📡</span> My Alert Zones</h3>
          <div class="centers-list">
            <div v-if="userData.alert_centers.length === 0" class="no-data">
              No saved zones. Click on the map to add one.
            </div>
            
            <div v-for="center in userData.alert_centers" :key="center.id" class="center-item">
              <div class="center-info">
                <span class="coords">📍 {{ center.lat.toFixed(4) }}, {{ center.lng.toFixed(4) }}</span>
                <div class="badges">
                  <span class="badge radius">{{ center.radius }} km</span>
                  <span class="badge magnitude" v-if="center.min_magnitude">M {{ center.min_magnitude }}</span>
                </div>
              </div>
              <button type="button" @click="removeCenter(center.id)" class="btn-delete-small" title="Delete zone">🗑</button>
            </div>
          </div>
        </div>

        <button type="submit" class="save-btn">Save Changes</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.account-wrapper {
  min-height: 100vh;
  overflow-y: auto;
  background-color: #1a1a2e;
  padding: 4rem 1rem; 
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
  position: relative;
}

.account-wrapper::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('/earthquake_photo_3.jpg');
  background-size: cover;
  background-position: center;
  background-attachment: scroll;
  filter: blur(1px);
  opacity: 0.4;
  z-index: 0;
}

.settings-card { 
  position: relative;
  z-index: 1;
  background: rgba(22, 33, 62, 0.95); 
  width: 100%; 
  max-width: 650px; 
  border-radius: 12px; 
  padding: 2.5rem; 
  border: 1px solid #2a3158; 
  color: #fff; 
  box-shadow: 0 10px 30px rgba(0,0,0,0.5); 
  margin-bottom: 4rem;
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
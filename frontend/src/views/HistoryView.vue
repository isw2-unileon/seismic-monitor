<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { apiService } from '../services/api';

const router = useRouter();
const earthquakes = ref([]);
const loading = ref(true);
const error = ref(null);

const minMagnitude = ref(1);
const hours = ref(1);

const fetchEarthquakes = async () => {
  try {
    loading.value = true;
    error.value = null;
    const data = await apiService.getEarthquakesHistory(minMagnitude.value, hours.value);
    if (data && data.features) {
      earthquakes.value = data.features;
    }
  } catch (err) {
    error.value = "Failed to load recent earthquakes.";
    console.error(err);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchEarthquakes);

const applyFilters = () => {
  fetchEarthquakes();
};

const goBack = () => router.push({ name: 'map' });

// Helper function to format the time
const formatTime = (isoString) => {
  const date = new Date(isoString);
  return date.toLocaleString();
};

</script>

<template>
  <div class="view-container">
    <header class="view-header">
      <button @click="goBack" class="back-btn">← Volver al Mapa</button>
      <h1>Historial Sísmico</h1>
    </header>

    <div class="filters-container">
      <div class="filter-item">
        <label for="magnitude-range">Magnitud Mínima: {{ minMagnitude.toFixed(1) }}</label>
        <input 
          id="magnitude-range" 
          type="range" 
          v-model.number="minMagnitude" 
          min="1" 
          max="10" 
          step="0.1"
        >
      </div>
      
      <div class="filter-item">
        <label for="time-range">Ventana de Tiempo: Últimas {{ hours }} {{ hours == 1 ? 'hora' : 'horas' }}</label>
        <input 
          id="time-range" 
          type="range" 
          v-model.number="hours" 
          min="1" 
          max="24" 
          step="1"
        >
      </div>

      <button @click="applyFilters" class="apply-btn">Aplicar Filtros</button>
    </div>
    
    <div v-if="loading" class="loading">
      Cargando terremotos recientes...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else-if="earthquakes.length === 0" class="no-data">
      No se han registrado terremotos con magnitud {{ minMagnitude }} o superior en las últimas {{ hours }} {{ hours == 1 ? 'hora' : 'horas' }}.
    </div>
    
    <div v-else class="earthquake-list">
      <div v-for="eq in earthquakes" :key="eq.properties.id" class="earthquake-card">
        <div class="eq-header">
          <span class="magnitude" :class="{'high-mag': eq.properties.magnitude >= 5}">
            M {{ eq.properties.magnitude.toFixed(1) }}
          </span>
          <span class="time">{{ formatTime(eq.properties.time) }}</span>
        </div>
        <div class="eq-details">
          <p><strong>Ubicación:</strong> {{ eq.properties.place }}</p>
          <p><strong>Profundidad:</strong> {{ eq.properties.depth_km || eq.geometry.coordinates[2] }} km</p>
          <p class="coords">
            <strong>Coordenadas:</strong> 
            {{ eq.geometry.coordinates[1].toFixed(4) }}, {{ eq.geometry.coordinates[0].toFixed(4) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.view-container {
  padding: 2rem;
  color: #e0e0e0;
  background-color: #1a1a2e;
  min-height: 100vh;
  font-family: system-ui, -apple-system, sans-serif;
}

.view-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
  border-bottom: 1px solid #333;
  padding-bottom: 1rem;
}

h1 {
  color: #fff;
  margin: 0;
}

.back-btn {
  background: transparent;
  border: 1px solid #e94560;
  color: #e94560;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: 0.3s;
  white-space: nowrap;
}

.back-btn:hover {
  background: #e94560;
  color: white;
}

.filters-container {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
  background-color: #16213e;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
  align-items: flex-end;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  border: 1px solid #2a3158;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
  min-width: 200px;
}

.filter-item label {
  font-weight: bold;
  color: #a0aab2;
}

input[type="range"] {
  width: 100%;
  cursor: pointer;
  accent-color: #e94560;
}

.apply-btn {
  background-color: #e94560;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s;
  height: fit-content;
}

.apply-btn:hover {
  background-color: #d63d56;
}

.apply-btn:active {
  background-color: #b02a42;
}

.loading, .error, .no-data {
  padding: 2rem;
  background-color: #16213e;
  border-radius: 8px;
  text-align: center;
  font-size: 1.1rem;
  border: 1px solid #2a3158;
}

.error {
  color: #ff6b6b;
  border: 1px solid #ff6b6b;
}

.earthquake-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 5rem;
}

.earthquake-card {
  background-color: #16213e;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s;
  border-left: 4px solid #e94560;
  border-top: 1px solid #2a3158;
  border-right: 1px solid #2a3158;
  border-bottom: 1px solid #2a3158;
}

.earthquake-card:hover {
  transform: translateY(-2px);
}

.eq-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  border-bottom: 1px solid #2a3158;
  padding-bottom: 0.75rem;
}

.magnitude {
  font-size: 1.25rem;
  font-weight: bold;
  color: #4cd137;
  background-color: rgba(76, 209, 55, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

.magnitude.high-mag {
  color: #e84118;
  background-color: rgba(232, 65, 24, 0.1);
}

.time {
  font-size: 0.9rem;
  color: #a0a0b0;
}

.eq-details p {
  margin: 0.5rem 0;
  line-height: 1.4;
}

.eq-details strong {
  color: #a0aab2;
}

.coords {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 0.75rem !important;
}
</style>
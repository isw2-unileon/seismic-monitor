import { ref } from 'vue'
import L from 'leaflet'

export function useAlertCenters() {
  const selectedMarkerId = ref(null)
  const pendingLocation = ref(null)

  const loadUserCenters = () => {
    const data = JSON.parse(localStorage.getItem('user_data') || '{}')
    return {
      centers: data.alert_centers || [],
      defaultRadius: data.alert_radius_km || 100
    }
  }

  const saveCenter = (latlng, radius) => {
    const { centers } = loadUserCenters()
    const newCenter = { id: Date.now(), lat: latlng.lat, lng: latlng.lng, radius }
    const data = JSON.parse(localStorage.getItem('user_data') || '{}')
    data.alert_centers = [...centers, newCenter]
    localStorage.setItem('user_data', JSON.stringify(data))
    return newCenter
  }

  return { selectedMarkerId, pendingLocation, loadUserCenters, saveCenter }
}
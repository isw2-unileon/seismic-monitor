// src/services/api.js

// Importación estática de los contratos JSON
import earthquakesMock from './mocks/earthquakes.json';

// Función auxiliar privada para simular latencia de red (ej. 800ms)
// Esto es vital para comprobar que tu interfaz muestra estados de "Cargando..."
const simulateNetworkLatency = (ms = 800) => new Promise(resolve => setTimeout(resolve, ms));

export const apiService = {
  /**
   * Obtiene la lista de sismos recientes en formato GeoJSON.
   */
  async getEarthquakes() {
    await simulateNetworkLatency();
    // Simula un código HTTP 200 OK devolviendo el mock directamente
    return earthquakesMock;
  },

  /**
   * Simula el inicio de sesión de un usuario o conecta con el backend.
   * @param {Object} credentials - { email, password }
   */
  async login(credentials) {
    if (!credentials.email || !credentials.password) {
      throw new Error("Missing email or password");
    }

    const response = await fetch('/api/v1/users/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || "Authentication failed");
    }

    return await response.json();
  },

  /**
   * Registers a new user.
   * @param {Object} credentials - { email, password }
   */
  async register(credentials) {
    if (!credentials.email || !credentials.password) {
      throw new Error("Missing email or password");
    }

    const response = await fetch('/api/v1/users/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || "Registration failed");
    }

    return await response.json();
  }
};
// src/services/api.js

// Importación estática de los contratos JSON
import earthquakesMock from './mocks/earthquakes.json';
import authMock from './mocks/auth.json';

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
   * Simula el inicio de sesión de un usuario.
   * @param {Object} credentials - { email, password }
   */
  async login(credentials) {
    await simulateNetworkLatency();
    
    // Validación defensiva básica simulada
    if (!credentials.email || !credentials.password) {
      throw new Error("Missing email or password"); // Simula un HTTP 400
    }

    // Devuelve el mock con el JWT y los datos del usuario
    return authMock;
  }
};
// src/services/api.js

// Importación estática de los contratos JSON
import earthquakesMock from './mocks/earthquakes.json';
import authMock from './mocks/auth.json';

const simulateNetworkLatency = (ms = 800) => new Promise(resolve => setTimeout(resolve, ms));
const API_BASE_URL = '/api/v1';

/**
 * Helper to transform static mock data into dynamic data with recent timestamps.
 */
const prepareMockData = (data) => {
  const now = new Date();
  return {
    ...data,
    features: data.features.map((feature, index) => ({
      ...feature,
      properties: {
        ...feature.properties,
        // Generates times: 5m ago, 15m ago, 25m ago...
        time: new Date(now.getTime() - (index * 10 + 5) * 60 * 1000).toISOString()
      }
    }))
  };
};

export const apiService = {
  async getEarthquakes() {
    try {
      const response = await fetch(`${API_BASE_URL}/earthquakes`);
      if (!response.ok) throw new Error('Network response was not ok');
      return await response.json();
    } catch (error) {
      console.warn("Using dynamic mock data due to API failure");
      await simulateNetworkLatency();
      return prepareMockData(earthquakesMock);
    }
  },

  async getEarthquakesHistory() {
    try {
      const response = await fetch(`${API_BASE_URL}/earthquakes/history`);
      if (!response.ok) throw new Error('Network response was not ok');
      return await response.json();
    } catch (error) {
      console.warn("Using dynamic mock data for history due to API failure");
      await simulateNetworkLatency();
      return prepareMockData(earthquakesMock);
    }
  },

  async login(credentials) {
    await simulateNetworkLatency();
    if (!credentials.email || !credentials.password) {
      throw new Error("Missing email or password");
    }
    return authMock;
  }
};

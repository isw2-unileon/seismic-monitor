// src/services/api.js

const API_BASE_URL = '/api/v1';

/**
 * Generates the USGS API URL for earthquakes with custom filters.
 */
const getUSGSUrl = (minMagnitude = 1, hours = 1) => {
  const startTime = new Date(Date.now() - hours * 60 * 60 * 1000).toISOString();
  return `https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&minmagnitude=${minMagnitude}&orderby=time&starttime=${startTime}`;
};

/**
 * Maps USGS GeoJSON data to the format expected by the frontend.
 */
const mapUSGSData = (data) => {
  if (!data || !data.features) return { type: 'FeatureCollection', features: [] };
  
  return {
    ...data,
    features: data.features.map(feature => ({
      ...feature,
      properties: {
        ...feature.properties,
        // Map USGS 'mag' to 'magnitude' expected by the frontend
        magnitude: feature.properties.mag,
        // Ensure time is in a format the frontend can handle (USGS returns ms)
        time: feature.properties.time,
        // Map top-level id to properties.id for component compatibility
        id: feature.id || feature.properties.code
      }
    }))
  };
};

export const apiService = {
  async getEarthquakes() {
    try {
      const response = await fetch(getUSGSUrl());
      if (!response.ok) throw new Error('USGS API response was not ok');
      const data = await response.json();
      return mapUSGSData(data);
    } catch (error) {
      console.error("Error fetching earthquakes from USGS:", error);
      return { type: 'FeatureCollection', features: [] };
    }
  },

  async getEarthquakesHistory(minMagnitude = 1, hours = 1) {
    try {
      const response = await fetch(getUSGSUrl(minMagnitude, hours));
      if (!response.ok) throw new Error('USGS API response was not ok');
      const data = await response.json();
      return mapUSGSData(data);
    } catch (error) {
      console.error("Error fetching earthquake history from USGS:", error);
      return { type: 'FeatureCollection', features: [] };
    }
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

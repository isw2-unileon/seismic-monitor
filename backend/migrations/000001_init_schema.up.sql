-- Allows the creation of the POSTGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    alert_radius_km NUMERIC NOT NULL DEFAULT 100,
    min_magnitude_alert NUMERIC(3,1) NOT NULL DEFAULT 3.0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create user's earthquake alert locations table
CREATE TABLE IF NOT EXISTS user_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(100), -- Opcional: "Casa", "Trabajo", etc.
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    alert_radius_km NUMERIC NOT NULL DEFAULT 50,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create earthquakes table with the USGS ID as primary key
CREATE TABLE IF NOT EXISTS earthquake (
    usgs_id VARCHAR(50) PRIMARY KEY,
    richter_scale NUMERIC(4,2) NOT NULL,
    place_name VARCHAR(255) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    depth_km NUMERIC NOT NULL,
    ocurred_at TIMESTAMPTZ NOT NULL,
    ingested_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create reported earthquakes table for user-reported events
CREATE TABLE reported_earthquakes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    reported_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 5. Create indexes
CREATE INDEX IF NOT EXISTS idx_user_locations_geom ON user_locations USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_earthquake_geom ON earthquake USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_earthquake_time ON earthquake (ocurred_at DESC);
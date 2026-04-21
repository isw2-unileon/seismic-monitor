-- Allows the creation of the POSTGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Create users table
CREATE TABLE users (
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

-- 3. Crete spacial index for users
CREATE INDEX idx_users_location ON users USING GIST (location);

-- 4. Create earthquakes table with the USGS ID as primary key
CREATE TABLE earthquake (
    usgs_id VARCHAR(50) PRIMARY KEY,
    richter_scale NUMERIC(4,2) NOT NULL,
    place_name VARCHAR(255) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    depth_km NUMERIC NOT NULL,
    ocurred_at TIMESTAMPTZ NOT NULL,
    ingested_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 5. Crear los índices para sismos
CREATE INDEX idx_earthquake_location ON earthquake USING GIST (location);
CREATE INDEX idx_earthquake_time ON earthquake (ocurred_at DESC);
-- migrations/002_create_location_tables.sql
-- +goose Up
CREATE TABLE locations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lat DOUBLE PRECISION NOT NULL CHECK (lat >= -90 AND lat <= 90),
  lng DOUBLE PRECISION NOT NULL CHECK (lng >= -180 AND lng <= 180),
  accuracy DOUBLE PRECISION NOT NULL CHECK (accuracy >= 0),

  recorded_at TIMESTAMPTZ NOT NULL,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  
  device_id UUID NOT NULL REFERENCES devices(id) ON DELETE SET NULL,
  person_id UUID REFERENCES people(id) ON DELETE SET NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX locations_person_recorded_at_idx
  ON locations (person_id, recorded_at DESC);

CREATE INDEX locations_device_recorded_at_idx
  ON locations (device_id, recorded_at DESC);

CREATE INDEX locations_received_at_idx
  ON locations (received_at DESC);

-- +goose Down
DROP TABLE IF EXISTS locations;

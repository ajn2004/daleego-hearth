-- migrations/001_create_people_devices_tables.sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE device_platform as ENUM (
  'android',
  'ios',
  'desktop',
  'server',
  'other'
);

CREATE TABLE people (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  display_name TEXT NOT NULL check (length(trim(display_name)) > 0),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  deactivated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ

);

CREATE TABLE devices (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  person_id UUID REFERENCES people(id) ON DELETE SET NULL,
  
  name TEXT NOT NULL,
  platform device_platform NOT NULL DEFAULT 'other',

  last_seen_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  deactivated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ

);

CREATE INDEX device_person_id_idx ON devices (person_id);
CREATE INDEX device_last_seen_at_idx ON devices (last_seen_at);

CREATE TABLE device_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,

    key_hash TEXT NOT NULL,
    key_prefix TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,

    UNIQUE (key_hash)
);

CREATE INDEX device_api_keys_device_id_idx ON device_api_keys (device_id);
CREATE INDEX device_api_keys_key_prefix_idx
  ON device_api_keys (key_prefix)
  WHERE revoked_at IS NULL;

create table device_pairing_codes (
  id UUID primary key default gen_random_uuid(),

  person_id UUID NOT NULL references people(id) on DELETE cascade,

  code_hash TEXT not null unique,

  used_by_device UUID references devices(id) on delete set null,

  expires_at TIMESTAMPTZ not null,
  used_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ not null default now()
);

CREATE INDEX device_pairing_codes_person_id_idx
  on device_pairing_codes (person_id);

CREATE index device_pairing_codes_expires_at_idx
  on device_pairing_codes (expires_at);

-- +goose Down
DROP table if exists device_pairing_codes;
DROP TABLE IF EXISTS device_api_keys;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS people;

DROP TYPE IF EXISTS device_platform;

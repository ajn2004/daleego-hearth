-- queries/locations.sql

-- name: CreateLocation :one
INSERT INTO locations (
  lat,
  lng,
  accuracy,
  recorded_at,
  device_id,
  person_id
) VALUES (
  sqlc.arg('lat'),
  sqlc.arg('lng'),
  sqlc.arg('accuracy'),
  sqlc.arg('recorded_at'),
  sqlc.arg('device_id')::uuid,
  sqlc.arg('person_id')::uuid
)
RETURNING *;

-- name: GetLocationByID :one
SELECT *
  FROM locations
 WHERE id = sqlc.arg('location_id')::uuid
   AND person_id = sqlc.arg('person_id')::uuid;

-- name: GetPersonLocations :many
SELECT *
  FROM locations
 WHERE person_id = sqlc.arg('person_id')::uuid;

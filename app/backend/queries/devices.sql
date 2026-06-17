-- queries/devices.sql

-- name: CreateDevice :one
INSERT INTO devices (
  name,
  platform,
  person_id
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('platform')::device_platform,
  sqlc.arg('person_id')::uuid
)
RETURNING *;

-- name: GetAllDevices :many
select * from devices;

-- name: GetDeviceByID :one
SELECT *
FROM devices
WHERE id = sqlc.arg('device_id')::uuid
  AND deleted_at IS NULL;

-- name: GetDevicesByPersonID :many
SELECT *
FROM devices
WHERE person_id = sqlc.arg('person_id')::uuid
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteDevice :one
UPDATE devices
SET
  deactivated_at = now(),
  deleted_at = now(),
  updated_at = now()
WHERE id = sqlc.arg('device_id')::uuid
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdateDeviceLastSeen :one
UPDATE devices
SET
  last_seen_at = now(),
  updated_at = now()
WHERE id = sqlc.arg('device_id')::uuid
  AND person_id = sqlc.arg('person_id')::uuid
  AND deleted_at IS NULL
      RETURNING *;

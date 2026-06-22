-- queries/device_api_keys.sql

-- name: CreateDeviceAPIKey :one
INSERT INTO device_api_keys (
  device_id,
  key_hash,
  key_prefix
) VALUES (
  sqlc.arg('device_id')::uuid,
  sqlc.arg('key_hash'),
  sqlc.arg('key_prefix')
)
RETURNING *;

-- name: GetActiveDeviceAPIKeysByPrefix :one
SELECT *
FROM device_api_keys
WHERE key_prefix = sqlc.arg('key_prefix')
  AND revoked_at IS NULL;

-- name: MarkDeviceAPIKeyUsed :exec
UPDATE device_api_keys
SET last_used_at = now()
WHERE id = sqlc.arg('device_api_key_id')::uuid;

-- name: RevokeDeviceAPIKeys :many
UPDATE device_api_keys
SET revoked_at = now()
WHERE device_id = sqlc.arg('device_id')::uuid
  AND revoked_at IS NULL
RETURNING *;

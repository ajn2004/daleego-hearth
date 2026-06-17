-- queries/device_pairing_codes.sql

-- name: CreateDevicePairingCode :one
INSERT INTO device_pairing_codes (
  person_id,
  code_hash,
  expires_at
) VALUES (
  sqlc.arg('person_id')::uuid,
  sqlc.arg('code_hash'),
  sqlc.arg('expires_at')
)
RETURNING *;

-- name: GetActivePairings :many
SELECT *
  FROM device_pairing_codes
 where used_at is null
   and expires_at > now();

-- name: GetValidDevicePairingCodeByHash :one
SELECT *
FROM device_pairing_codes
WHERE code_hash = sqlc.arg('code_hash')
  AND used_at IS NULL
  AND expires_at > now();

-- name: GetPairingCodeByID :one
SELECT *
  FROM device_pairing_codes
 where id = sqlc.arg('pairing_code_id')::uuid;

-- name: MarkDevicePairingCodeUsed :one
UPDATE device_pairing_codes
SET
  used_at = now(),
  used_by_device_id = sqlc.arg('device_id')::uuid
WHERE id = sqlc.arg('pairing_code_id')::uuid
  AND used_at IS NULL
RETURNING *;

-- name: RevokePairingCode :one
UPDATE device_pairing_codes
   set expires_at = now()
 where used_at is null
   and id = sqlc.arg('pairing_code_id')::uuid
       RETURNING *;

-- name: GetExpiredPairings :many
SELECT *
  FROM device_pairing_codes
 where expires_at < now();

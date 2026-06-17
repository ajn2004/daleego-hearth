-- queries/people.sql
-- name: CreatePerson :one
INSERT INTO people (
  display_name
) VALUES (
  sqlc.arg('display_name')
)
RETURNING *;

-- name: GetPersonByID :one
SELECT *
  FROM people
 WHERE id = sqlc.arg('person_id')::uuid
   AND deleted_at IS NULL;

-- name: UpdatePersonName :one
UPDATE people
   SET
display_name = sqlc.arg('display_name'),
updated_at = now()
 WHERE id = sqlc.arg('person_id')::uuid
   AND deleted_at IS NULL
       RETURNING *;

-- name: SetPersonToDeleted :one
UPDATE people
   SET
deleted_at = now(),
updated_at = now()
 WHERE id = sqlc.arg('person_id')
   AND deleted_at IS NULL
       RETURNING *;

-- name: AdminDeletePerson :exec
DELETE FROM people WHERE id = sqlc.arg('person_id')::uuid;

-- name: AdminGetAllPeople :many
SELECT *
  FROM people
 ORDER BY created_at DESC;

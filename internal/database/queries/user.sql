-- name: CreateUser :one
INSERT INTO users(
    user_email,
    user_password,
    user_status,
    user_role
) VALUES(
    $1,$2,$3,$4
) RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET
    user_password = COALESCE(sqlc.narg(user_password), user_password),
    user_status = COALESCE(sqlc.narg(user_status), user_status),
    user_role = COALESCE(sqlc.narg(user_role), user_role)
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :one
UPDATE users
SET
    user_deleted_at = now()
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NULL
RETURNING *;

-- name: CleanSoftDelete :one
DELETE FROM users
WHERE
    user_deleted_at IS NOT NULL AND user_uuid = sqlc.arg(user_uuid)::uuid
RETURNING *;

-- name: RestoreUser :one
UPDATE users
SET
    user_deleted_at = NULL
WHERE
    user_deleted_at IS NOT NULL AND user_uuid = sqlc.arg(user_uuid)::uuid
RETURNING *;

-- name: GetUserByUUID :one
SELECT *
FROM users
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid
    AND user_deleted_at IS NULL;
-- name: CountRecords :one
SELECT count(*)
FROM users
WHERE (
    sqlc.narg(deleted)::bool IS NULL
    OR (user_deleted_at IS NOT NULL AND sqlc.narg(deleted)::bool IS TRUE)
    OR (user_deleted_at IS NULL AND sqlc.narg(deleted)::bool IS FALSE)
) AND (
    sqlc.narg(search)::TEXT IS NULL
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_fullname ILIKE '%' || sqlc.narg(search) || '%');

-- name: CreateUserToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at
)
VALUES ($1, NOW(), NOW(), $2, NOW() + INTERVAL '60 days')
RETURNING *;

-- name: GetUserByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeUserToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;
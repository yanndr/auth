-- name: GetUser :one
SELECT *
FROM users
WHERE username = ?
LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, password_hash)
VALUES (?, ?);
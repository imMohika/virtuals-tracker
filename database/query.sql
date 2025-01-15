-- name: GetAgentByID :one
SELECT *
FROM agents
WHERE id = ?
LIMIT 1;

-- name: GetAgentByUID :one
SELECT *
FROM agents
WHERE uid = ?
LIMIT 1;

-- name: ExistsAgentByUID :one
SELECT COUNT(1)
FROM agents
WHERE uid = ?;

-- name: CreateAgent :one
INSERT INTO agents (uid, name, status, category, mcap, notified)
VALUES (?,?,?,?,?,?)
RETURNING *;

-- name: UpdateAgent :exec
UPDATE agents
set notified = ?
WHERE id = ?;

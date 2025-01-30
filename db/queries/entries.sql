-- name: CreateEntry :one
INSERT INTO entries(account_id, amount)
VALUES ($1, $2)
returning *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: DeleteAllEntries :exec
DELETE FROM entries;
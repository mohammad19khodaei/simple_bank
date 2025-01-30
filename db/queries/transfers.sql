-- name: CreateTransfer :one
INSERT INTO transfers(from_account_id,to_account_id,amount)
VALUES ($1,$2,$3)
returning *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: DeleteAllTransfers :exec
DELETE FROM transfers;
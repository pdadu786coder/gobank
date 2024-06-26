-- name: CreateTransfer :one
INSERT INTO transfers(
  from_account_id, to_account_id, amount
)VALUES(
    $1,$2,$3
)
RETURNING *;

-- name: GetTransfer :one
select * from transfers
where id = $1 limit 1;

-- name: ListTransfers :many
select * from transfers
where  
    from_account_id = $1 OR
    to_account_id = $2
ORDER by id
LIMIT $3
OFFSET $4;
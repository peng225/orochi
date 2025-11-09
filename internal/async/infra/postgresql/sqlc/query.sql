-- name: SelectBucket :one
SELECT * FROM bucket
WHERE id = $1;

-- name: DeleteBucket :exec
DELETE FROM bucket
WHERE id = $1;

-- name: SelectJobs :many
SELECT * FROM job
LIMIT $1;

-- name: DeleteJob :exec
DELETE FROM job
WHERE id = $1;

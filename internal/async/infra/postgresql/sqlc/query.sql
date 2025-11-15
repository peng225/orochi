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

-- name: SelectDatastores :many
SELECT * FROM datastore;

-- name: UpdateDatastoreStatus :exec
UPDATE datastore
SET status = $1
WHERE id = $2;

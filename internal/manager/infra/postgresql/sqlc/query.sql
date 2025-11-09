-- name: InsertDatastore :one
INSERT INTO datastore (
   base_url
) VALUES (
  $1
)
RETURNING id;

-- name: SelectDatastore :one
SELECT * FROM datastore
WHERE id = $1;

-- name: SelectDatastoreIDs :many
SELECT id FROM datastore;

-- name: InsertLocationGroup :one
INSERT INTO location_group (
   current_datastores,
   desired_datastores
) VALUES (
  $1, $2
)
RETURNING id;

-- name: UpdateDesiredDatastores :exec
UPDATE location_group
SET desired_datastores = $1
WHERE id = $2;

-- name: SelectLocationGroups :many
SELECT * from location_group;

-- name: InsertBucket :one
INSERT INTO bucket (
   name, status
) VALUES (
  $1, 'active'
)
RETURNING id;

-- name: UpdateBucketStatus :exec
UPDATE bucket
SET status = $1
WHERE id = $2;

-- name: SelectBucket :one
SELECT * FROM bucket
WHERE id = $1;

-- name: InsertJob :one
INSERT INTO job (
   kind, data
) VALUES (
  $1, $2
)
RETURNING id;

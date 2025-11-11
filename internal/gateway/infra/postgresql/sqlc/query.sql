-- name: SelectDatastores :many
SELECT * FROM datastore;

-- name: CreateObjectMetadata :one
INSERT INTO object_metadata (
   name,
   bucket_id,
   location_group_id
) VALUES (
  $1, $2, $3
)
RETURNING id;

-- name: SelectObjectMetadataByName :many
SELECT * FROM object_metadata
WHERE name = $1 AND bucket_id = $2;

-- name: DeleteObjectMetadata :exec
DELETE FROM object_metadata
WHERE id = $1;

-- name: SelectObjectMetadatas :many
SELECT * FROM object_metadata
WHERE id >= $1 AND bucket_id = $2
LIMIT $3;

-- name: SelectLocationGroupsByECConfigID :many
SELECT * from location_group
WHERE ec_config_id = $1;

-- name: SelectLocationGroup :one
SELECT * from location_group
WHERE id = $1;

-- name: SelectBucketByName :one
SELECT * FROM bucket
WHERE name = $1;

-- name: SelectECConfig :one
SELECT * FROM ec_config
WHERE id = $1;

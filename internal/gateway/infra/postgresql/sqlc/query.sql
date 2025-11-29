-- name: InsertObjectMetadata :one
INSERT INTO object_metadata (
   name,
   status,
   bucket_name,
   location_group_id
) VALUES (
  $1, 'creating', $2, $3
)
RETURNING id;

-- name: SelectObjectMetadataByName :many
SELECT * FROM object_metadata
WHERE name = $1 AND bucket_name = $2;

-- name: UpdateObjectMetadataStatus :exec
UPDATE object_metadata
SET status = $1
WHERE id = $2;

-- name: DeleteObjectMetadata :exec
DELETE FROM object_metadata
WHERE id = $1;

-- name: SelectObjectMetadatas :many
SELECT * FROM object_metadata
WHERE id >= $1 AND bucket_name = $2
LIMIT $3;

-- name: InsertObjectVersion :one
INSERT INTO object_version (
   update_time,
   object_id
) VALUES (
  $1, $2
)
RETURNING id;

-- name: DeleteObjectVersionsByObjectID :exec
DELETE FROM object_version
WHERE object_id = $1;

-- name: SelectLocationGroupsByECConfigID :many
SELECT * from location_group
WHERE ec_config_id = $1;

-- name: SelectLocationGroup :one
SELECT * from location_group
WHERE id = $1;

-- name: SelectECConfig :one
SELECT * FROM ec_config
WHERE id = $1;

-- name: SelectECConfigByNumbers :one
SELECT * FROM ec_config
WHERE num_data = $1 AND num_parity = $2;

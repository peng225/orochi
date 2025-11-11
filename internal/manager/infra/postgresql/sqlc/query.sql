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
   desired_datastores,
   ec_config_id
) VALUES (
  $1, $2, $3
)
RETURNING id;

-- name: UpdateDesiredDatastores :exec
UPDATE location_group
SET desired_datastores = $1
WHERE id = $2;

-- name: SelectLocationGroupsByECConfigID :many
SELECT * from location_group
WHERE ec_config_id = $1;

-- name: InsertBucket :one
INSERT INTO bucket (
   name, ec_config_id, status
) VALUES (
  $1, $2, 'active'
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

-- name: InsertECConfig :one
INSERT INTO ec_config (
   num_data,
   num_parity
) VALUES (
  $1, $2
)
RETURNING id;

-- name: SelectECConfigByNumbers :one
SELECT * FROM ec_config
WHERE num_data = $1 AND num_parity = $2;

-- name: SelectECConfigs :many
SELECT * from ec_config;

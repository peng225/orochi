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

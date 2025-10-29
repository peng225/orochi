CREATE TABLE datastore(
    id BIGSERIAL PRIMARY KEY,
    base_url VARCHAR(128) NOT NULL
);

CREATE TABLE location_group(
    id BIGSERIAL PRIMARY KEY,
    current_datastores BIGINT[] NOT NULL,
    desired_datastores BIGINT[] NOT NULL
);

CREATE TABLE bucket(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE
);

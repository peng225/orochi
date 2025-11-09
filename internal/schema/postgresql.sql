CREATE TABLE datastore(
    id BIGSERIAL PRIMARY KEY,
    base_url VARCHAR(128) NOT NULL
);

CREATE TABLE location_group(
    id BIGSERIAL PRIMARY KEY,
    current_datastores BIGINT[] NOT NULL,
    desired_datastores BIGINT[] NOT NULL
);

CREATE TYPE bucket_status AS ENUM ('active', 'deleted');
CREATE TABLE bucket(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    status bucket_status NOT NULL
);

CREATE TABLE object_metadata(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    bucket_id BIGINT NOT NULL,
    location_group_id BIGINT NOT NULL,
    FOREIGN KEY (bucket_id) REFERENCES bucket(id),
    FOREIGN KEY (location_group_id) REFERENCES location_group(id)
);

CREATE TABLE job(
    id BIGSERIAL PRIMARY KEY,
    kind VARCHAR(128) NOT NULL,
    data JSONB NOT NULL
);

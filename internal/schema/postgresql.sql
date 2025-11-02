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

CREATE TABLE object_metadata(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    bucket_id BIGINT NOT NULL,
    location_group_id BIGINT NOT NULL,
    FOREIGN KEY (bucket_id) REFERENCES bucket(id),
    FOREIGN KEY (location_group_id) REFERENCES location_group(id)
);

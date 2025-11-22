CREATE TYPE datastore_status AS ENUM ('active', 'down');
CREATE TABLE datastore(
    id BIGSERIAL PRIMARY KEY,
    base_url VARCHAR(128) NOT NULL UNIQUE,
    status datastore_status NOT NULL
);

CREATE TABLE ec_config(
    id BIGSERIAL PRIMARY KEY,
    num_data INTEGER NOT NULL,
    num_parity INTEGER NOT NULL,
    UNIQUE(num_data, num_parity)
);


CREATE TYPE location_group_status AS ENUM ('active', 'deleting');
CREATE TABLE location_group(
    id BIGSERIAL PRIMARY KEY,
    datastores BIGINT[] NOT NULL,
    ec_config_id BIGINT NOT NULL,
    status location_group_status NOT NULL,
    FOREIGN KEY (ec_config_id) REFERENCES ec_config(id)
);

CREATE TYPE bucket_status AS ENUM ('active', 'deleted');
CREATE TABLE bucket(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    ec_config_id BIGINT NOT NULL,
    status bucket_status NOT NULL,
    FOREIGN KEY (ec_config_id) REFERENCES ec_config(id)
);

CREATE TYPE object_status AS ENUM ('creating', 'active');
CREATE TABLE object_metadata(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    status object_status NOT NULL,
    bucket_id BIGINT NOT NULL,
    location_group_id BIGINT NOT NULL,
    FOREIGN KEY (bucket_id) REFERENCES bucket(id),
    FOREIGN KEY (location_group_id) REFERENCES location_group(id)
);

CREATE TABLE object_version(
    id BIGSERIAL PRIMARY KEY,
    update_time TIMESTAMP NOT NULL,
    object_id BIGINT NOT NULL,
    FOREIGN KEY (object_id) REFERENCES object_metadata(id)
);

CREATE TABLE job(
    id BIGSERIAL PRIMARY KEY,
    kind VARCHAR(128) NOT NULL,
    data JSONB NOT NULL
);

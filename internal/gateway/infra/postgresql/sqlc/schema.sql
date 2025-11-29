CREATE TABLE IF NOT EXISTS ec_config(
    id BIGSERIAL PRIMARY KEY,
    num_data INTEGER NOT NULL,
    num_parity INTEGER NOT NULL,
    UNIQUE(num_data, num_parity)
);

CREATE TYPE location_group_status AS ENUM ('active', 'deleting');
CREATE TABLE IF NOT EXISTS location_group(
    id BIGSERIAL PRIMARY KEY,
    datastores BIGINT[] NOT NULL,
    ec_config_id BIGINT NOT NULL,
    status location_group_status NOT NULL,
    FOREIGN KEY (ec_config_id) REFERENCES ec_config(id)
);

CREATE TYPE object_status AS ENUM ('creating', 'updating', 'active');
CREATE TABLE IF NOT EXISTS object_metadata(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    status object_status NOT NULL,
    bucket_name VARCHAR(128) NOT NULL,
    location_group_id BIGINT NOT NULL,
    FOREIGN KEY (location_group_id) REFERENCES location_group(id)
);

CREATE TABLE IF NOT EXISTS object_version(
    id BIGSERIAL PRIMARY KEY,
    update_time TIMESTAMP NOT NULL,
    object_id BIGINT NOT NULL,
    FOREIGN KEY (object_id) REFERENCES object_metadata(id)
);

CREATE TABLE IF NOT EXISTS job(
    id BIGSERIAL PRIMARY KEY,
    kind VARCHAR(128) NOT NULL,
    data JSONB NOT NULL
);

CREATE TYPE datastore_status AS ENUM ('active', 'down');
CREATE TABLE IF NOT EXISTS datastore(
    id BIGSERIAL PRIMARY KEY,
    base_url VARCHAR(128) NOT NULL UNIQUE,
    status datastore_status NOT NULL
);

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

CREATE TYPE bucket_status AS ENUM ('active', 'deleted');
CREATE TABLE IF NOT EXISTS bucket(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    ec_config_id BIGINT NOT NULL,
    status bucket_status NOT NULL,
    FOREIGN KEY (ec_config_id) REFERENCES ec_config(id)
);

CREATE TABLE IF NOT EXISTS job(
    id BIGSERIAL PRIMARY KEY,
    kind VARCHAR(128) NOT NULL,
    data JSONB NOT NULL
);

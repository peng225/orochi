create table datastore(
    id bigserial primary key,
    base_url varchar(128) not null
);

create table location_group(
    id bigserial primary key,
    current_datastores bigint[] not null,
    desired_datastores bigint[] not null
)

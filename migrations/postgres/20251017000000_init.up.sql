create schema if not exists app;

create table if not exists app.users
(
    id         serial primary key,
    name       text      not null,
    email      text      not null,
    password   text      not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);

alter table app.users
    add constraint users_email_key unique (email);

create table if not exists app.tables
(
    id          text                     not null primary key,
    name        text                     not null,
    database_id integer                  not null,
    columns     jsonb                    not null default '[]',
    created_at  timestamp with time zone not null default now(),
    deleted_at  timestamp with time zone
);

create index if not exists tables_database_idx on app.tables (database_id);

create table if not exists app.users_databases
(
    user_id     integer                  not null,
    database_id integer                  not null,
    role        text                     not null,
    created_at  timestamp with time zone not null default now(),
    deleted_at  timestamp with time zone
);

alter table app.users_databases
    add constraint users_databases_pkey primary key (user_id, database_id);

create table if not exists app.databases
(
    id         serial primary key,
    name       text                     not null,
    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);

create schema if not exists users_tablespace;
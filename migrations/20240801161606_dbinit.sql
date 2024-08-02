-- +goose Up
CREATE TABLE IF NOT EXISTS
    users (
        id varchar(26) primary key,
        username text not null unique,
        password text,
        created_at timestamptz default current_timestamp,
        updated_at timestamptz default current_timestamp
    );

CREATE TABLE IF NOT EXISTS
    todos (
        id varchar(26) primary key,
        description text not null,
        is_timed boolean default false,
        deadline timestamptz,
        user_id varchar(26) not null references users (id) on delete set null,
        created_at timestamptz default current_timestamp,
        updated_at timestamptz default current_timestamp
    );

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
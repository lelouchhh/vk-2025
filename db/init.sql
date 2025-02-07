CREATE TABLE IF NOT EXISTS containers (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(255) NOT NULL,
    last_ping TIMESTAMP NOT NULL DEFAULT NOW(),
    ping_time FLOAT not null,
    status BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS account (
    id serial primary key,
    login varchar(255) not null,
    password varchar(255) not null
);
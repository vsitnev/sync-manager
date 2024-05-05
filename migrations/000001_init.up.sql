CREATE SCHEMA IF NOT EXISTS md;
CREATE SCHEMA IF NOT EXISTS exchange;

CREATE TABLE IF NOT EXISTS md."user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    "password" VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS exchange."message" (
    id BIGSERIAL PRIMARY KEY,
    routing VARCHAR(128) NOT NULL,
    "message" jsonb NOT NULL,
    dead BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
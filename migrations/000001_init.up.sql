CREATE SCHEMA IF NOT EXISTS md;
CREATE SCHEMA IF NOT EXISTS exchange;

CREATE TABLE IF NOT EXISTS md."user"
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(128) NOT NULL,
    "password" VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS exchange."message"
(
    id         BIGSERIAL PRIMARY KEY,
    routing    VARCHAR(128) NOT NULL,
    "message"  jsonb        NOT NULL,
    dead       BOOLEAN      NOT NULL DEFAULT false,
    retried    BOOLEAN      NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS exchange."source"
(
    id             BIGSERIAL PRIMARY KEY,
    "name"         VARCHAR(128) NOT NULL,
    description    VARCHAR(256) NOT NULL,
    code           VARCHAR(8)   NOT NULL,
    receive_method varchar(16)  NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_code ON exchange."source" (code);

CREATE TABLE IF NOT EXISTS exchange."route"
(
    id        BIGSERIAL PRIMARY KEY,
    "name"    VARCHAR(128) NOT NULL,
    url       VARCHAR(256) NOT NULL,
    source_fk int8         NOT NULL,
    CONSTRAINT source_fk FOREIGN KEY (source_fk) REFERENCES exchange.source (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_name ON exchange."route" (name);

-- Timestamps are TIMESTAMPTZ and stored in UTC; time zones are a
-- presentation concern handled by the client.
CREATE TABLE users (
    id         UUID        PRIMARY KEY DEFAULT uuidv7(),
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

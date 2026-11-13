ALTER DATABASE SET udt_postgres SET timezone TO 'America/São_Paulo';

CREATE TABLE users (
    id         uuid default uuidv7() PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

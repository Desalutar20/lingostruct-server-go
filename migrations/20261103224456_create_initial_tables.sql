-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    email TEXT NOT NULL UNIQUE,
    first_name TEXT,
    last_name TEXT,
    hashed_password TEXT,
    role TEXT NOT NULL DEFAULT 'regular' CHECK (role in ('admin', 'regular')),
    is_banned BOOLEAN NOT NULL DEFAULT false,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    avatar_url TEXT,
    google_id TEXT UNIQUE,
    github_id text UNIQUE
);

CREATE TABLE IF NOT EXISTS outbox (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    processed_at TIMESTAMPTZ,
    data JSONB NOT NULL,
    type TEXT NOT NULL CHECK (type in ('email'))
);

CREATE TABLE IF NOT EXISTS workspace (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    street TEXT NOT NULL,
    street_number TEXT NOT NULL,
    postal_code TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS workspace;

DROP TABLE IF EXISTS outbox;

DROP TABLE IF EXISTS users;

-- +goose Up
CREATE TYPE provider AS ENUM ('credential', 'google', 'github');

CREATE TABLE "users" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  image TEXT,
  is_banned BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE "accounts" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL REFERENCES "users" (id) ON DELETE CASCADE,
  account_id TEXT NOT NULL,
  provider_id provider NOT NULL,
  password TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT account_provider_unique UNIQUE (provider_id, account_id),
  CONSTRAINT password_required_when_credentials CHECK (
    provider_id <> 'credential' OR password IS NOT NULL
  )
);

CREATE TABLE "verifications" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL REFERENCES "users" (id) ON DELETE CASCADE,
  identifier TEXT NOT NULL,
  value TEXT NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT verification_identifier_value_unique UNIQUE (identifier, value)
);

CREATE TABLE "sessions" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL REFERENCES "users" (id) ON DELETE CASCADE,
  ip_address TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "verifications";
DROP TABLE IF EXISTS "accounts";
DROP TABLE IF EXISTS "users";

DROP TYPE IF EXISTS provider;


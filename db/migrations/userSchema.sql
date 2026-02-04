CREATE TABLE users (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  password_hash Text Not Null,
  email      TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
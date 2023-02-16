CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE Table if not exists post(
  id serial PRIMARY KEY,
  customer_id UUID,
  name VARCHAR(200),
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
CREATE Table if NOT exists media(
  id serial PRIMARY KEY,
  post_id int REFERENCES post(id),
  name VARCHAR(100),
  link TEXT,
  type VARCHAR(100)
);
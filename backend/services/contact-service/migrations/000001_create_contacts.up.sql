CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS contacts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  adder_id TEXT NOT NULL,
  added_id TEXT NOT NULL,
  alias_name VARCHAR(255),
  created_at TIMESTAMP DEFAULT now(),

  UNIQUE (adder_id, added_id)
);

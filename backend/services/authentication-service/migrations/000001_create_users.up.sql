CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email TEXT UNIQUE NOT NULL,
  image_url TEXT,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE users(
  id BIGSERIAL PRIMARY KEY,
  login TEXT,
  password VARCHAR
);

CREATE TABLE product(
  id BIGSERIAL PRIMARY KEY,
  name TEXT
);

CREATE TABLE producr_image(
  id BIGSERIAL PRIMARY KEY,
  product_id INTEGER REFERENCES product (id) ON DELETE CASCADE,
  file_name VARCHAR(255) NOT NULL,
  original_name VARCHAR(255),
  s3_bucket VARCHAR(255) NOT NULL,
  s3_key VARCHAR(500) NOT NULL,
  mime_type VARCHAR(100),
  file_size BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE product(
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  description TEXT,
  parameters TEXT,
  count INTEGER CHECK (count >= 0) NOT NULL,
  price INTEGER CHECK (count >= 0) NOT NULL
);

CREATE TABLE product_image(
  id BIGSERIAL PRIMARY KEY,
  product_id INTEGER REFERENCES product (id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  key VARCHAR(500) NOT NULL UNIQUE
);

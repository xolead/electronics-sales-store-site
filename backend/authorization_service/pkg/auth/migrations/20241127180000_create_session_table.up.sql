-- +goose Up
-- +goose StatementBegin

CREATE TABLE refresh(
  id BIGSERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  refresh_token VARCHAR,
  expires_at TIMESTAMP,
  worker BOOLEAN
);

CREATE TABLE session(
  id BIGSERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  active BOOLEAN
);

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employees
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(64)         NOT NULL,
    last_name  VARCHAR(64)         NOT NULL,
    email      VARCHAR(128) UNIQUE NOT NULL,
    hire_date  DATE                NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE socket_clients(
    id VARCHAR(255) PRIMARY KEY,
    data TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS socket_clients;
-- +goose StatementEnd

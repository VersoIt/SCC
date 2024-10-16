-- +goose Up
-- +goose StatementBegin
CREATE TABLE sockets(
    id VARCHAR(255) PRIMARY KEY,
    data BYTEA
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sockets;
-- +goose StatementEnd

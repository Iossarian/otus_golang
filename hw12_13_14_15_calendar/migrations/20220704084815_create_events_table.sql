-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id UUID NOT NULL,
    user_id int NOT NULL,
    title VARCHAR NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd

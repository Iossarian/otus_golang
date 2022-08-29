-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id                UUID      NOT NULL,
    user_id           int       NOT NULL,
    title             VARCHAR   NOT NULL,
    description       TEXT,
    start_date        TIMESTAMP NOT NULL,
    end_date          TIMESTAMP NOT NULL,
    notification_date TIMESTAMP NULL,
    is_notified       SMALLINT DEFAULT 0
);
CREATE UNIQUE INDEX idx_user_start_date
    ON events (user_id, start_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd

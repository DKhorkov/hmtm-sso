-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id                 SERIAL PRIMARY KEY,
    display_name       VARCHAR      NOT NULL,
    email              VARCHAR(255) NOT NULL UNIQUE,
    email_confirmed    BOOLEAN      NOT NULL DEFAULT FALSE,
    password           VARCHAR(255) NOT NULL,
    phone              VARCHAR(50),
    phone_confirmed    BOOLEAN      NOT NULL DEFAULT FALSE,
    telegram           VARCHAR(70) UNIQUE,
    telegram_confirmed BOOLEAN      NOT NULL DEFAULT FALSE,
    avatar             VARCHAR UNIQUE,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

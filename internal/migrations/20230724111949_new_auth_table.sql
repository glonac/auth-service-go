-- миграции точно не должны лежать в internal
-- https://github.com/golang-standards/project-layout
-- путь в Makefile не актуален

-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE auths
(
        id              bigserial
        primary key,
        created_at      timestamp with time zone,
        updated_at      timestamp with time zone,
        deleted_at      timestamp with time zone,
        password        text,
        email           text,
        code            text,
        code_created_at timestamp with time zone,
        is_verified     boolean,
        update_at       timestamp with time zone
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
-- нет скриптов отката

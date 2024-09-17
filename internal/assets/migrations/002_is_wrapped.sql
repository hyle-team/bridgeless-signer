-- +migrate Up

ALTER TABLE deposits
    ADD COLUMN is_wrapped_token boolean DEFAULT false;

-- +migrate Down

ALTER TABLE deposits
    DROP COLUMN is_wrapped_token;
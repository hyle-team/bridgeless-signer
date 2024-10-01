-- +migrate Up

ALTER TABLE deposits
    ADD COLUMN signature text;

-- +migrate Down

ALTER TABLE deposits
    DROP COLUMN signature;
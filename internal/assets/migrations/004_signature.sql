-- +migrate Up

ALTER TABLE deposits
    ADD COLUMN signature text, ADD COLUMN tx_nonce int;

-- +migrate Down

ALTER TABLE deposits
    DROP COLUMN signature, DROP COLUMN tx_nonce;
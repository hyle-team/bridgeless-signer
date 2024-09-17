-- +migrate Up

ALTER TABLE deposits
    RENAME COLUMN amount TO deposit_amount;

ALTER TABLE deposits
    ADD COLUMN withdrawal_amount TEXT;

-- +migrate Down

ALTER TABLE deposits
    RENAME COLUMN deposit_amount to amount;

ALTER TABLE deposits
    DROP COLUMN withdrawal_amount;
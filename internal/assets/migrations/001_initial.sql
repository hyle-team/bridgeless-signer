-- +migrate Up

CREATE TABLE deposits
(
    id                  BIGSERIAL PRIMARY KEY,
    tx_hash             VARCHAR(66) not null,
    tx_event_id         int         not null,
    chain_id            text        not null,
    status              int         not null,
    withdrawal_tx_hash  VARCHAR(66),
    withdrawal_chain_id text,
    CONSTRAINT unique_deposit UNIQUE (tx_hash, tx_event_id, chain_id)
)

-- +migrate Down

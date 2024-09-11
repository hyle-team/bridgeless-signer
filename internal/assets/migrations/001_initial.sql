-- +migrate Up

CREATE TABLE deposits
(
    id                  BIGSERIAL PRIMARY KEY,
    tx_hash             text not null,
    tx_event_id         int         not null,
    chain_id            text        not null,
    status              int         not null,

    depositor          text,
    amount             text,
    deposit_token      text,
    receiver           text,
    withdrawal_token   text,
    deposit_block      BIGINT,

    withdrawal_tx_hash  text,
    withdrawal_chain_id text,

    is_wrapped_token    bool,
    submit_status       int         not null,

    CONSTRAINT unique_deposit UNIQUE (tx_hash, tx_event_id, chain_id)
);

-- +migrate Down

drop table deposits;

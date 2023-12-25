create table transaction
(
    id                    serial primary key,
    currency              varchar(3)     not null,
    amount                numeric(15, 2) not null,
    wallet_or_card_number varchar(255)   not null,
    status                varchar(15)    not null,
    created_at            timestamp default current_timestamp
);

create table balance
(
    id             serial primary key,
    currency       varchar(3)     not null,
    balance        numeric(15, 2) not null,
    frozen_balance numeric(15, 2) not null default 0,
    updated_at     timestamp               default current_timestamp
);

CREATE INDEX idx_transactions_currency ON transaction (currency);
CREATE INDEX idx_balances_currency ON balance (currency);
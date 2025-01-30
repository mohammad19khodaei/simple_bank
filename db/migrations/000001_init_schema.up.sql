-- CREATE DATABASE simple_bank;

CREATE TABLE accounts(
    id serial PRIMARY KEY,
    owner varchar NOT NULL,
    balance bigint NOT NULL,
    currency varchar NOT NULL ,
    created_at timestamptz default now()
);

CREATE TABLE entries(
    id serial PRIMARY KEY,
    account_id int NOT NULL ,
    amount bigint NOT NULL,
    created_at timestamptz default now(),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE transfers(
    id serial PRIMARY KEY ,
    from_account_id int NOT NULL ,
    to_account_id int NOT NULL ,
    amount bigint NOT NULL ,
    created_at timestamptz default now(),
    FOREIGN KEY (from_account_id) REFERENCES accounts(id),
    FOREIGN KEY (to_account_id) REFERENCES accounts(id)
);

CREATE INDEX ON accounts (owner);
CREATE INDEX ON entries (account_id);
CREATE INDEX ON transfers(from_account_id);
CREATE INDEX ON transfers(to_account_id);
CREATE INDEX ON transfers(from_account_id,to_account_id);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';
COMMENT ON COLUMN transfers.amount IS 'must be positive';
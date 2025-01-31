CREATE TABLE users(
    username varchar PRIMARY KEY,
    hashed_password varchar NOT NULL,
    full_name varchar NOT NULL,
    email varchar UNIQUE NOT NULL,
    password_changed_at timestamptz default now(),
    created_at timestamptz default now()
);

ALTER TABLE accounts ADD FOREIGN KEY (owner) REFERENCES users(username);

ALTER TABLE accounts ADD CONSTRAINT "owner_currency_key" UNIQUE (owner, currency);
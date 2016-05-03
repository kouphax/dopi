CREATE TABLE users(
    id                    SERIAL PRIMARY KEY,
    name                  TEXT UNIQUE,
    email                 TEXT UNIQUE,
    password              TEXT,
    confirm_token         TEXT,
    confirmed             BOOLEAN,
    attempt_number        INT,
    attempt_time          TIMESTAMP,
    locked                TIMESTAMP,
    recover_token         TEXT,
    recover_token_expiry  TIMESTAMP
)
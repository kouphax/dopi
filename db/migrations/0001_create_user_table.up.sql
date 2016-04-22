CREATE TABLE users(
    id                  SERIAL PRIMARY KEY,
    name                TEXT UNIQUE,
    email               TEXT UNIQUE,
    password            TEXT,
    confirmToken        TEXT,
    confirmed           BOOLEAN,
    attemptNumber       INT,
    attemptTime         TIMESTAMP,
    locked              TIMESTAMP,
    recoverToken        TEXT,
    recoverTokenExpiry  TIMESTAMP
)
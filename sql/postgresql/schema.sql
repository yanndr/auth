CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    username      text NOT NULL UNIQUE CHECK (username <> ''),
    password_hash text NOT NULL CHECK (password_hash <> '')
);

CREATE INDEX username_idx ON users (username);

CREATE TABLE version
(
    version text NOT NULL DEFAULT '0.0.0'
);

INSERT into version
VALUES ('0.1');
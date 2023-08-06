CREATE TABLE users
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    username      text NOT NULL,
    password_hash text NOT NULL,
    UNIQUE(username)
);

CREATE INDEX username_idx ON users (username);

CREATE TABLE version
(
    version text NOT NULL DEFAULT '0.0.0'
);

INSERT into version
VALUES ('0.1');
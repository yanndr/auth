CREATE TABLE users
(
    id       BIGSERIAL PRIMARY KEY,
    username text NOT NULL UNIQUE ,
    password_hash text NOT NULL
);

CREATE TABLE version
(
    version text NOT NULL DEFAULT '0.0.0'
);

INSERT into version VALUES ('0.1');
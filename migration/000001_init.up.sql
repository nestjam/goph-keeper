BEGIN;

CREATE TABLE users(
    user_id      UUID PRIMARY KEY       DEFAULT gen_random_uuid(),
    email        VARCHAR(64) UNIQUE     NOT NULL CHECK ( email <> '' ),
    password     VARCHAR(250)           NOT NULL CHECK ( octet_length(password) <> 0 )
);

CREATE TABLE keys(
    key_id                  UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    key_data                BYTEA,
    encriptions_count       INTEGER,
    encrypted_data_size     BIGINT,
    is_disposed             BOOLEAN                 DEFAULT 'false'
);

CREATE TABLE secrets(
    secret_id       UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    user_id         UUID                    NOT NULL REFERENCES users (user_id),
    key_id          UUID                    NOT NULL REFERENCES keys (key_id),
    name            TEXT,
    data            BYTEA
);

END;
DROP TABLE IF EXISTS expressions;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    login    TEXT,
    password TEXT
);

CREATE TABLE expressions
(
    id                   SERIAL PRIMARY KEY,
    value                TEXT,
    answer               FLOAT,
    logs                 TEXT,
    ready                INT,
    alive_expires_at     BIGINT,
    creation_time        TEXT,
    end_calculation_time TEXT,
    server_name          TEXT,
    user_id              INT,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);

CREATE TABLE operations
(
    id            SERIAL PRIMARY KEY,
    time_add      INT,
    time_subtract INT,
    time_divide   INT,
    time_multiply INT,
    user_id       INT,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);
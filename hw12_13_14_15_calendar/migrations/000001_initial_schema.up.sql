CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id   UUID PRIMARY KEY,
    name VARCHAR(256)
);

CREATE TABLE events
(
    id          UUID PRIMARY KEY,
    title       TEXT      NOT NULL,
    start_dt    TIMESTAMP NOT NULL,
    end_dt      TIMESTAMP NOT NULL,
    description TEXT,
    user_id     UUID      NOT NULL,
    notify INTERVAL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);
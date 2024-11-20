CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE events
(
    id          UUID PRIMARY KEY,
    title       TEXT      NOT NULL,
    start_dt    TIMESTAMP NOT NULL,
    end_dt      TIMESTAMP NOT NULL,
    description TEXT,
    user_id     UUID      NOT NULL,
    notify_at INTERVAL
);


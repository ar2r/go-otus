CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id   UUID PRIMARY KEY,
    name VARCHAR(256)
);
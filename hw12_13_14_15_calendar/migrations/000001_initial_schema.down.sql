-- Disable foreign key constraints
SET CONSTRAINTS ALL DEFERRED;

-- Drop tables
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS events;
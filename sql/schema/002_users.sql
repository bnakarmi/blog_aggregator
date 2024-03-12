-- +goose Up
ALTER TABLE USERS
    ADD COLUMN API_KEY VARCHAR(64) UNIQUE NOT NULL;

UPDATE USERS 
SET API_KEY = ENCODE(SHA256(RANDOM()::TEXT::BYTEA), 'HEX');

-- +goose Down
ALTER TABLE USERS
    DROP COLUMN API_KEY


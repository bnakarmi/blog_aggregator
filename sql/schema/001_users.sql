-- +goose Up
CREATE TABLE USERS (
    ID UUID PRIMARY KEY,
    CREATED_AT TIMESTAMP NOT NULL,
    UPDATED_AT TIMESTAMP NOT NULL,
    NAME VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE USERS;

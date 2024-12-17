-- +goose Up
CREATE TABLE students
(
    id         SERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- +goose Down
DROP TABLE students;

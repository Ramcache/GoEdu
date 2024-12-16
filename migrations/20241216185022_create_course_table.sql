-- +goose Up
CREATE TABLE courses
(
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL
);


-- +goose Down
DROP TABLE courses;
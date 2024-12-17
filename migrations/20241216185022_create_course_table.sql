-- +goose Up
CREATE TABLE courses
(
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    description   TEXT NOT NULL,
    instructor_id INT NOT NULL,
    CONSTRAINT fk_instructor
        FOREIGN KEY (instructor_id)
            REFERENCES instructors(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE courses;

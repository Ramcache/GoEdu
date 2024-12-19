-- +goose Up
CREATE TABLE courses
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    description   TEXT NOT NULL,
    instructor_id INT NOT NULL,
    CONSTRAINT fk_instructor
        FOREIGN KEY (instructor_id)
            REFERENCES instructors(id)
            ON DELETE CASCADE,
    CONSTRAINT courses_name_not_empty CHECK (char_length(name) > 0),
    CONSTRAINT courses_description_not_empty CHECK (char_length(description) > 0)
);

CREATE UNIQUE INDEX courses_name_unique ON courses (name);

-- +goose Down
DROP TABLE courses;

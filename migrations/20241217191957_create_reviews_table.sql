-- +goose Up
CREATE TABLE reviews (
                         id          SERIAL PRIMARY KEY,
                         student_id  INT NOT NULL,
                         course_id   INT NOT NULL,
                         comment     TEXT,
                         rating      INT CHECK (rating >= 1 AND rating <= 5),
                         created_at  TIMESTAMP DEFAULT NOW(),
                         FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
                         FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE reviews;

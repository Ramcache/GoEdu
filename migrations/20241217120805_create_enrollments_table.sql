-- +goose Up
CREATE TABLE enrollments (
                             student_id INT NOT NULL,
                             course_id INT NOT NULL,
                             enrolled_at TIMESTAMP DEFAULT NOW(),
                             PRIMARY KEY (student_id, course_id),
                             FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
                             FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE enrollments;

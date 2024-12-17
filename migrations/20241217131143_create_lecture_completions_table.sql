-- +goose Up
CREATE TABLE lecture_completions
(
    student_id   INT NOT NULL,
    lecture_id   INT NOT NULL,
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (student_id, lecture_id),
    FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
    FOREIGN KEY (lecture_id) REFERENCES lectures (id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE lecture_completions;

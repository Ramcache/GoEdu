-- +goose Up
CREATE TABLE lectures (
                          id SERIAL PRIMARY KEY,
                          course_id INT NOT NULL,
                          title TEXT NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW(),
                          FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE lectures;

-- +goose Up
CREATE TABLE students (
                          id SERIAL PRIMARY KEY,           -- Уникальный идентификатор студента
                          name TEXT NOT NULL,              -- Имя студента
                          email TEXT UNIQUE NOT NULL,      -- Email студента (должен быть уникальным)
                          password TEXT NOT NULL,          -- Хэш пароля студента
                          created_at TIMESTAMP DEFAULT NOW() -- Дата создания записи
);


-- +goose Down
DROP TABLE students;

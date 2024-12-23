## GoEdu Project 🚀

GoEdu — это образовательный проект, разработанный на языке Go с использованием gRPC и соблюдением принципов SOLID. Проект включает:

- **gRPC API** для взаимодействия с сервисами 💻.
- **SOLID-ориентированную структуру** для поддержки масштабируемости и чистоты кода 🛠️.
- **Поддержку PostgreSQL** с использованием библиотеки `pgx` 🐘.
- Логирование через **zap** 📝.
- **Swagger-документация** для gRPC API 📑.

---

## Структура проекта 🗂️

```
# Структура проекта

<details>
  <summary>Главные каталоги</summary>
  - **cmd/**: Точка входа в приложение
  - **internal/**: Основная логика и бизнес-слой
  - **logs/**: Лог-файлы приложения
  - **migrations/**: Миграции базы данных
  - **proto/**: gRPC протоколы
  - **swagger/**: Документация API с использованием Swagger
</details>

<details>
  <summary>Подробнее о каталогах и файлах</summary>

  ### cmd/server/
  - `main.go`: Точка входа в серверное приложение

  ### internal/config/
  - `config.go`: Конфигурация приложения

  ### internal/logger/
  - `logger.go`: Настройки логирования через zap

  ### internal/models/
  - `course.go`: Модель для курсов
  - `Instructor.go`: Модель для преподавателей
  - `lecture.go`: Модель для лекций
  - `review.go`: Модель для отзывов
  - `students.go`: Модель для студентов

  ### internal/repository/
  - `course_repository.go`: Репозиторий для курсов
  - `database.go`: Управление подключениями к базе данных
  - `enrollments_repository.go`: Репозиторий для регистраций студентов
  - `instructor_repository.go`: Репозиторий для преподавателей
  - `lecture_repositry.go`: Репозиторий для лекций
  - `review_repository.go`: Репозиторий для отзывов
  - `student_repository.go`: Репозиторий для студентов

  ### internal/service/
  - `education_service.go`: Сервис для работы с курсами
  - `enrollments_service.go`: Сервис для работы с регистрациями студентов
  - `instructor_service.go`: Сервис для преподавателей
  - `lecture_service.go`: Сервис для лекций
  - `review_service.go`: Сервис для отзывов
  - `student_service.go`: Сервис для студентов

  ### migrations/
  - `20241216185022_create_course_table.sql`: Миграция для создания таблицы курсов
  - `20241217060156_create_students_table.sql`: Миграция для создания таблицы студентов

  ### proto/
  - `education.proto`: Определение gRPC-сообщений
  - `education.pb.go`: Сгенерированный Go-код для сообщений gRPC
  - `education_grpc.pb.go`: Сгенерированный Go-код для сервисов gRPC

  ### swagger/
  - `education.swagger.json`: JSON описание API
  - `swagger-ui/`: Интерфейс Swagger UI для документации
</details>

```

---

## Установка и запуск 🚀

1. **Клонирование репозитория** 📥
   ```bash
   git clone https://github.com/Ramcache/GoEdu.git
   cd GoEdu
   ```

2. **Настройка окружения** ⚙️
    - Создайте файл `.env` на основе примера и укажите настройки для базы данных и сервера.

   Пример `.env`:
   ```
   DATABASE_URL=postgres://name:password@localhost:5432/dbname
   GRPC_PORT=50051
   JWT_SECRET_KEY=Your_Secret_key
   TOKEN_EXPIRATION_HOURS=24
   HTTP_PORT=8080
   ```

3. **Установка зависимостей** 📦
   ```bash
   go mod tidy
   ```

4. **Генерация gRPC-кода** 🛠️
   Используйте `protoc` для генерации gRPC кода:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/education.proto
   ```

5. **Запуск сервера** 🚀
   ```bash
   go run cmd/server/main.go
   ```

---

## Логирование с использованием zap 📝

Проект использует библиотеку `zap` для эффективного и структурированного логирования. Конфигурация логирования настраивается в файле конфигурации и используется в каждом сервисе для записи ключевых событий.

Пример использования в сервисе:
```go
import (
    "go.uber.org/zap"
)

func (s *EducationService) CreateCourse(ctx context.Context, req *CreateCourseRequest) (*CreateCourseResponse, error) {
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    logger.Info("Creating new course", zap.String("course_name", req.CourseName))
    
    // Логика создания курса
}
```

---

## Swagger-документация 📑

Для удобства использования gRPC API проект включает интеграцию с **Swagger**, что позволяет взаимодействовать с API через веб-интерфейс. Для доступа к документации откройте браузер и перейдите по следующему URL:

```
http://localhost:8080/swagger-ui/
```

---

## Makefile для автоматизации задач ⚙️

Проект включает `Makefile` для автоматизации различных задач, таких как генерация gRPC кода и запуск сервера. Пример содержимого Makefile:

```Makefile
# Makefile for GoEdu project

.PHONY: generate run

generate:
    protoc --go_out=. --go-grpc_out=. proto/education.proto

run:
    go run cmd/server/main.go
```

---

## Принципы разработки 📏

Проект построен с использованием **SOLID-принципов**:
- **S** — Single Responsibility Principle: каждый сервис отвечает за свою область (студенты, курсы и т.д.) 🎯.
- **O** — Open/Closed Principle: сервисы легко расширяются без изменения базового кода 🔧.
- **L** — Liskov Substitution Principle: интерфейсы позволяют подмену реализаций 🔄.
- **I** — Interface Segregation Principle: интерфейсы минимизированы и конкретны 🎛️.
- **D** — Dependency Inversion Principle: зависимости внедряются через интерфейсы 🔌.

---

## gRPC API 🎯

Файл `proto/education.proto` описывает структуру запросов и ответов. Пример:

```proto
service EducationService {
    rpc GetCourse (GetCourseRequest) returns (GetCourseResponse);
    rpc CreateCourse (CreateCourseRequest) returns (CreateCourseResponse);
}

message GetCourseRequest {
    string course_id = 1;
}

message GetCourseResponse {
    string course_id = 1;
    string course_name = 2;
}
```

---

## TODO 📝
- [ ] Улучшить производительность запросов к базе данных 📊.
- [ ] Реализовать дополнительные gRPC методы 🏆.

---

## Авторы 🤝
- **Ramzan**

---

**GoEdu** — обучайся и развивайся с лёгкостью! 🎓✨

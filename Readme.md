## GoEdu Project 🚀

GoEdu — это образовательный проект, разработанный на языке Go с использованием gRPC и соблюдением принципов SOLID. Проект включает:

- **gRPC API** для взаимодействия с сервисами 💻.
- **SOLID-ориентированную структуру** для поддержки масштабируемости и чистоты кода 🛠️.
- **Поддержку PostgreSQL** с использованием библиотеки `pgx` 🐘.
- Логирование через **zap** 📝.
- **Swagger-документация** для gRPC API 📑.
- **Docker-контейнеры** для унифицированного развертывания и работы с проектом 🐳.

---

## Структура проекта 🗂️

```
GoEdu/
├── cmd/                           # Точка входа в приложение
│   └── server/                    # Содержит главный серверный файл
│       └── main.go                # Точка входа в серверное приложение
├── internal/                      # Основная логика и бизнес-слой
│   ├── config/                    # Конфигурационные файлы
│   │   └── config.go              # Конфигурация приложения (например, параметры подключения к БД)
│   ├── logger/                    # Логирование
│   │   └── logger.go              # Настройки и реализация логирования через zap
│   ├── middleware/                # Промежуточное ПО для обработки запросов
│   │   └── auth_middleware.go     # Мидлвар для авторизации пользователей
│   ├── models/                    # Модели данных, которые используются в приложении
│   │   ├── course.go              # Модель для курсов
│   │   ├── Instructor.go          # Модель для преподавателей
│   │   ├── lecture.go             # Модель для лекций
│   │   ├── review.go              # Модель для отзывов
│   │   └── students.go            # Модель для студентов
│   ├── repository/                # Репозитории для взаимодействия с базой данных
│   │   ├── course_repository.go   # Репозиторий для работы с курсами
│   │   ├── database.go            # Управление подключениями к базе данных
│   │   ├── enrollments_repository.go # Репозиторий для регистраций студентов
│   │   ├── instructor_repository.go  # Репозиторий для преподавателей
│   │   ├── lecture_repositry.go   # Репозиторий для лекций
│   │   ├── review_repository.go   # Репозиторий для отзывов
│   │   └── student_repository.go  # Репозиторий для студентов
│   └── service/                   # Сервисы, которые обрабатывают бизнес-логику
│       ├── education_service.go   # Сервис для работы с курсами
│       ├── educations_test.go     # Тесты для сервиса курсов
│       ├── enrollments_service.go # Сервис для работы с регистрациями студентов
│       ├── enrollments_test.go    # Тесты для сервиса регистраций
│       ├── instructor_service.go  # Сервис для работы с преподавателями
│       ├── instructor_test.go     # Тесты для сервиса преподавателей
│       ├── lecture_service.go     # Сервис для работы с лекциями
│       ├── lecture_test.go        # Тесты для сервиса лекций
│       ├── main_test.go           # Основные тесты
│       ├── review_service.go      # Сервис для работы с отзывами
│       ├── review_test.go         # Тесты для сервиса отзывов
│       ├── student_service.go     # Сервис для работы со студентами
│       └── student_test.go        # Тесты для сервиса студентов
├── migrations/                    # Миграции базы данных
│   ├── 20241216185022_create_course_table.sql  # Миграция для создания таблицы курсов
│   ├── 20241217060156_create_students_table.sql # Миграция для создания таблицы студентов
│   ├── 20241217120805_create_enrollments_table.sql # Миграция для регистрации студентов
│   ├── 20241217123415_create_lectures_table.sql  # Миграция для создания таблицы лекций
│   ├── 20241217131143_create_lecture_completions_table.sql # Миграция для завершенных лекций
│   ├── 20241217135242_create_instructors_table.sql # Миграция для создания таблицы преподавателей
│   └── 20241217191957_create_reviews_table.sql   # Миграция для создания таблицы отзывов
├── proto/                         # gRPC протоколы и файлы для генерации кода
│   ├── education.pb.go            # Сгенерированный Go-код для сообщений gRPC
│   ├── education.pb.gw.go         # Генерация кода для работы с API Gateway
│   ├── education.proto            # Определение gRPC-сообщений для курса
│   └── education_grpc.pb.go       # Сгенерированный Go-код для сервисов gRPC
|   └── swagger/                   # Документация API с использованием Swagger
│   ├── education.swagger.json     # JSON-файл описания API с использованием Swagger
│   ├── swagger-ui/                # Папка с файлами для отображения интерфейса Swagger UI
│   │   ├── favicon-16x16.png      # Иконка для интерфейса Swagger
│   │   ├── favicon-32x32.png      # Иконка для интерфейса Swagger
│   │   ├── index.css              # Стили для интерфейса Swagger
│   │   ├── index.html             # Главная HTML-страница интерфейса Swagger
│   │   ├── oauth2-redirect.html   # HTML-страница для редиректа OAuth2
│   │   ├── swagger-initializer.js # Скрипт инициализации Swagger UI
│   │   ├── swagger-ui.css         # Стили для Swagger UI
│   │   ├── swagger-ui.js          # Скрипт для отображения Swagger UI
│   │   ├── swagger-ui-bundle.js   # Скомпилированный скрипт для Swagger UI
│   │   ├── swagger-ui-es-bundle.js # Скомпилированный скрипт для европейской версии Swagger UI
│   │   ├── swagger-ui-es-bundle-core.js # Основной скрипт для европейской версии Swagger UI
│   │   ├── swagger-ui-standalone-preset.js # Скрипт для автономного использования Swagger UI
├── coverage.out                   # Файл покрытия тестами
├── go.mod                         # Модульные зависимости проекта
├── Makefile                       # Файл для автоматизации сборки и задач
├── Dockerfile                     # Файл для создания Docker-образа приложения
├── docker-compose.yml             # Файл для управления контейнерами с помощью Docker Compose
└── Readme.md                      # Документация проекта
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

5. **Сборка и запуск через Docker** 🐳
   Сборка и запуск контейнеров:
   ```bash
   docker-compose up --build
   ```

   После успешного запуска:
   - HTTP сервер будет доступен по адресу [http://localhost:8080](http://localhost:8080).
   - База данных PostgreSQL будет работать на порту `5432`.

6. **Локальный запуск без Docker** 🚀
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

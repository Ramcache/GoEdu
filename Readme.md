# GoEdu Project 🚀

GoEdu — это образовательный проект, разработанный на языке Go с использованием gRPC и соблюдением принципов SOLID. Проект включает:

- **gRPC API** для взаимодействия с сервисами 💻.
- **SOLID-ориентированную структуру** для поддержки масштабируемости и чистоты кода 🛠️.
- **Поддержку PostgreSQL** с использованием библиотеки `pgx` 🐘.
- Логирование через **zap** 📝.

---

## Структура проекта 🗂️

```
GoEdu/
|── cmd/
|   ── server/
|       ── main.go          # Точка входа сервера 🚀
|── proto/
|   ── education.proto    # Определение gRPC-сообщений и сервисов 📡
|── internal/
|   ── service/
|       ── education_service.go      # Сервис для управления курсами 🎓
|       ── enrollments_service.go    # Управление регистрациями студентов 🌟
|       ── instructor_service.go     # Сервис для инструкторов 👨‍👩‍👦
|       ── lecture_service.go        # Управление лекциями 📅
|       ── student_service.go        # Сервис для студентов 📚
|       ── review_service.go         # Отзывы студентов 📊
|── config/
|   ── config.go          # Конфигурация приложения ⚙️
|── .env                  # Переменные окружения 🌐
|── go.mod                # Зависимости проекта 📦
|── go.sum                # Контрольная сумма зависимостей ✅
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
- [ ] Добавить интеграционные тесты ✅.
- [ ] Настроить CI/CD для деплоя 🚀.
- [ ] Добавить документацию для API 📄.
- [ ] Улучшить производительность запросов к базе данных 📊.
- [ ] Реализовать дополнительные gRPC методы 🏆.
- [ ] Добавить логирование 🔗.
- [ ] Реализовать мониторинг 📊.
---

## Авторы 🤝
- **Ramzan**

---

**GoEdu** — обучайся и развивайся с лёгкостью! 🎓✨

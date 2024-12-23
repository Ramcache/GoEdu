# Makefile для автоматизации задач в проекте GoEdu

# Путь к директории с proto-файлами и миграциями
PROTO_DIR=proto
MIGRATIONS_DIR=migrations
BIN_DIR=bin

# Команда для protoc
PROTOC=protoc

# Команда для утилиты миграций
MIGRATE=migrate

# Пакет для скомпилированных Go-файлов
GO_OUT=$(PROTO_DIR)

# Цель по умолчанию
all: proto migrate

# Компиляция .proto файлов
proto:
	@echo "Компилируем .proto файлы..."
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		--swagger_out=$(PROTO_DIR)/swagger \
		$(PROTO_DIR)/*.proto
	@echo "Генерация .proto завершена."

# Применение миграций
migrate:
	@echo "Применяем миграции..."
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "postgres://<user>:<password>@<host>:<port>/<db_name>?sslmode=disable" up
	@echo "Миграции применены."

# Проверка тестов с покрытием
test:
	@echo "Запуск тестов..."
	go test ./... -cover
	@echo "Тестирование завершено."

# Очистка сгенерированных файлов
clean:
	@echo "Очистка сгенерированных файлов..."
	rm -rf $(PROTO_DIR)/*.pb.go $(PROTO_DIR)/*.gw.go $(PROTO_DIR)/swagger/*.json
	@echo "Очистка завершена."

# Запуск приложения
run:
	@echo "Запуск приложения..."
	go run cmd/server/main.go

.PHONY: all proto migrate test clean run

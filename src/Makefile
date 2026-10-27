GO          ?= go
SWAG        ?= swag
APP         ?= shopAPI
MAIN        ?= ./cmd/app/main.go
AUTH        ?= ./cmd/auth/main.go
PROTO_PATH  ?= proto/auth.proto
READONLY_DB ?= postgres://readonly_user:password@localhost:5432/shopapi?sslmode=disable

.PHONY: help run run-auth run2 run3 build proto tools swag clean

# Цель по умолчанию: справка
help: ## Отобразить список доступных команд
	@echo "Использование: make [цель]"
	@echo ""
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

swag: ## Перегенерировать Swagger-документацию (docs/)
	@echo "==> Генерация Swagger-документации..."
	$(SWAG) init -g $(MAIN)

run: swag ## Запустить магазин (нужен src/.env и запущенный PostgreSQL)
	@echo "==> Запуск: http://localhost:8080/swagger/index.html"
	$(GO) run $(MAIN)

run-auth: ## Запустить gRPC-сервис авторизации (:50051, нужна БД auth_db)
	@echo "==> gRPC-сервис авторизации: :$(shell grep -oP '(?<=AUTH_GRPC_PORT=)\d+' .env 2>/dev/null || echo 50051)"
	$(GO) run $(AUTH)

run2: ## Второй инстанс на :8081 (read-only доступ к БД) — для балансировки
	@echo "==> Инстанс :8081 (read-only)"
	PORT=8081 DATABASE_URL='$(READONLY_DB)' $(GO) run $(MAIN)

run3: ## Третий инстанс на :8082 (read-only доступ к БД) — для балансировки
	@echo "==> Инстанс :8082 (read-only)"
	PORT=8082 DATABASE_URL='$(READONLY_DB)' $(GO) run $(MAIN)

proto: ## Сгенерировать Go-код из proto/auth.proto
	@echo "==> Генерация pb-кода..."
	protoc --go_out=. --go_opt=module=shopAPI \
	       --go-grpc_out=. --go-grpc_opt=module=shopAPI \
	       $(PROTO_PATH)

tools: ## Установить инструменты: swag, protoc-gen-go, protoc-gen-go-grpc
	@echo "==> Установка инструментов..."
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build: swag ## Собрать бинарник shopAPI
	@echo "==> Сборка..."
	$(GO) build -o $(APP) $(MAIN)

clean: ## Удалить собранный бинарник
	@echo "==> Очистка..."
	rm -f $(APP)

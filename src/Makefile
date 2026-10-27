GO          ?= go
SWAG        ?= swag
APP         ?= shopAPI
MAIN        ?= ./cmd/app/main.go
READONLY_DB ?= postgres://readonly_user:password@localhost:5432/shopapi?sslmode=disable

.PHONY: help run run2 run3 build swag clean

# Цель по умолчанию: справка
help: ## Отобразить список доступных команд
	@echo "Использование: make [цель]"
	@echo ""
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

swag: ## Перегенерировать Swagger-документацию (docs/)
	@echo "==> Генерация Swagger-документации..."
	$(SWAG) init -g $(MAIN)

run: swag ## Запустить приложение (нужен src/.env и запущенный PostgreSQL)
	@echo "==> Запуск: http://localhost:8080/swagger/index.html"
	$(GO) run $(MAIN)

run2: ## Второй инстанс на :8081 (read-only доступ к БД) — для балансировки
	@echo "==> Инстанс :8081 (read-only)"
	PORT=8081 DATABASE_URL='$(READONLY_DB)' $(GO) run $(MAIN)

run3: ## Третий инстанс на :8082 (read-only доступ к БД) — для балансировки
	@echo "==> Инстанс :8082 (read-only)"
	PORT=8082 DATABASE_URL='$(READONLY_DB)' $(GO) run $(MAIN)

build: swag ## Собрать бинарник shopAPI
	@echo "==> Сборка..."
	$(GO) build -o $(APP) $(MAIN)

clean: ## Удалить собранный бинарник
	@echo "==> Очистка..."
	rm -f $(APP)

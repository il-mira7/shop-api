GO      ?= go
SWAG    ?= swag
APP     ?= shopAPI
MAIN    ?= ./cmd/app/main.go

.PHONY: help run build swag clean

help: ## Отобразить список доступных команд
	@echo "Использование: make [цель]"
	@echo ""
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

swag: ## Перегенерировать Swagger-документацию (docs/)
	@echo "==> Генерация Swagger-документации..."
	$(SWAG) init -g $(MAIN)

run: swag ## Запустить приложение (нужен src/.env и запущенный PostgreSQL)
	@echo "==> Запуск: http://localhost:8080/swagger/index.html"
	$(GO) run $(MAIN)

build: swag ## Собрать бинарник shopAPI
	@echo "==> Сборка..."
	$(GO) build -o $(APP) $(MAIN)

clean: ## Удалить собранный бинарник
	@echo "==> Очистка..."
	rm -f $(APP)

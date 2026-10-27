# ShopAPI

REST API магазина бытовой техники: клиенты, товары, поставщики и изображения товаров. Учебный проект [School 21](https://21-school.ru) (направление backend), выполнен самостоятельно.

Слоистая архитектура с чётким разделением контрактов: HTTP-хендлеры работают с DTO, сервисный слой — с доменными моделями, репозитории — с DAO; между слоями — мапперы. Схема БД описана миграциями, которые накатываются автоматически при старте; спецификация OpenAPI генерируется из swagger-аннотаций.

## Стек

- **Go 1.24**, фреймворк **Gin**
- **PostgreSQL** через `pgx/v5`
- **uber-fx** — сборка зависимостей
- **golang-migrate** — миграции (встроены в бинарник, применяются при старте)
- **Swagger/OpenAPI 3** — аннотации + Swagger UI
- `nyaruka/phonenumbers` — валидация телефонов поставщиков

## Запуск

Понадобятся Go 1.24+, запущенный PostgreSQL и установленный [swag](https://github.com/swaggo/swag):

```bash
createdb shopapi
cp .env.example src/.env         # и впишите свой DATABASE_URL
cd src
make run
```

- Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- База: миграции применяются автоматически при первом старте

## Публикация за Nginx (этап 2)

Конфигурация — [`src/nginx.conf`](src/nginx.conf). Nginx слушает `:80` (HTTP) и `:443` (HTTPS, самоподписанный сертификат локального домена `nginx.test`):

- маршрутизация `/api` → `/api/v1`, по `/api/v1/` отдаётся Swagger UI приложения;
- GET-запросы `/api/v1/` балансируются между тремя инстансами в соотношении **2:1:1** — `:8080` основной, `:8081`/`:8082` работают с read-only доступом к БД; мутации всегда уходят на основной инстанс (заголовок `X-Backend` показывает, кто отвечал);
- GET-ответы кэшируются: зона 10 МБ в памяти + 100 МБ на диске, TTL 10 минут, заголовок `X-Cache-Status`;
- `/` — статика из `var/www/static` (`expires 1h`), `/admin/` — pgAdmin, `/status` — страница статуса Nginx;
- gzip-сжатие ответов (текст, JSON, XML — без медиа-типов).

Запуск стенда:

```bash
# 1. read-only пользователь БД для инстансов :8081/:8082
sudo -u postgres psql -c "CREATE USER readonly_user WITH PASSWORD 'password';"
sudo -u postgres psql -c "GRANT CONNECT ON DATABASE shopapi TO readonly_user; GRANT USAGE ON SCHEMA public TO readonly_user; GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;"

# 2. три инстанса приложения (каждый в своём терминале)
make run     # :8080 — основной
make run2    # :8081 — read-only
make run3    # :8082 — read-only

# 3. nginx
sudo nginx -t -c $PWD/nginx.conf   # проверить конфигурацию
sudo nginx -c $PWD/nginx.conf      # запустить
sudo nginx -s reload               # перечитать после изменений
```

Проверить балансировку: `for i in {1..8}; do curl -s -o /dev/null -D - http://localhost/api/v1/products | grep -i x-backend; done`

## API

Все пути начинаются с префикса `/api/v1`.

### Клиенты

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/clients` | создать клиента |
| `GET` | `/clients` | все клиенты (опциональная пагинация `limit`/`offset`) |
| `GET` | `/clients/search?name=&surname=` | поиск по имени и фамилии |
| `PATCH` | `/clients/:id/address` | сменить адрес клиента |
| `DELETE` | `/clients/:id` | удалить клиента |

### Товары

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/products` | создать товар |
| `GET` | `/products` | все доступные товары |
| `GET` | `/products/:id` | товар по id |
| `PATCH` | `/products/:id/stock` | уменьшить остаток на складе |
| `DELETE` | `/products/:id` | удалить товар |

### Поставщики

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/suppliers` | создать поставщика |
| `GET` | `/suppliers` | все поставщики |
| `GET` | `/suppliers/:id` | поставщик по id |
| `PATCH` | `/suppliers/:id/address` | сменить адрес поставщика |
| `DELETE` | `/suppliers/:id` | удалить поставщика |

### Изображения товаров

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/images/product/:product_id` | загрузить изображение товара |
| `PUT` | `/images/:id` | заменить изображение |
| `GET` | `/images/product/:product_id` | изображение конкретного товара |
| `GET` | `/images/:id` | изображение по id |
| `DELETE` | `/images/:id` | удалить изображение |

Изображения отдаются массивом байт с заголовком `application/octet-stream` и загружаются браузером как файл.

Ошибки: невалидное тело запроса — `400`, сущность не найдена — `404`, пустой список — `200` с пустым массивом.

## Архитектура

```
src/
├── cmd/app/                 — точка входа, загрузка .env
├── nginx.conf               — конфигурация Nginx (этап 2)
├── var/www/static/          — статика для раздачи через Nginx
└── internal/
    ├── config/              — конфигурация (DATABASE_URL)
    ├── di/                  — сборка зависимостей (uber-fx)
    ├── domain/              — доменные модели и ошибки
    ├── handler/             — HTTP-хендлеры, DTO, маршрутизация
    ├── repository/postgres/ — DAO, репозитории, SQL-миграции
    └── service/             — бизнес-логика
```

## Что я сделала бы иначе

- DI на uber-fx для графа из одной зависимости (БД) — тяжеловато; хватило бы явной ручной сборки в `main.go`.
- Юнит-тестов нет — добавила бы табличные тесты на валидацию и сервисный слой, интеграционные — на репозитории против тестовой БД.

## Этапы развития

- ✅ **Этап 1**  — **REST API**: реализация API магазина на Go и PostgreSQL — CRUD четырёх сущностей, слои DTO/service/DAO, автомиграции, Swagger UI
  
  Требования задания: RESTful-методология, DTO и мапперы между слоями, паттерн Repository, полное покрытие методов спецификацией OpenAPI.
- ✅ **Этап 2**  — **Nginx**: сервис публикуется за обратным прокси на `:80`/`:443`
  
  Требования задания: маршрутизация `/api` → `/api/v1` → Swagger, раздача статики с `/` (index.html + image.png), `/admin` → pgAdmin, `/status` → страница статуса Nginx; балансировка GET-запросов на три инстанса приложения 2:1:1 (дополнительные — с read-only доступом к БД); кэширование GET-запросов; gzip-сжатие без медиа-типов; HTTPS с самоподписанным сертификатом (openssl + локальный DNS).
- 🔜 **Этап 3** — **Авторизация и DI**: аутентификация, JWT, RBAC, gRPC

Проект в работе — этапы будут дополняться.

---

*Учебный проект School 21.*


# WebPet1 — Task Manager API

Простой REST API для управления пользователями и их задачами (todo-list с привязкой задач к пользователю). Проект написан на Go в слоистой архитектуре (layered architecture) с генерацией серверного кода из OpenAPI-спецификации.

## Содержание

- [Описание проекта](#описание-проекта)
- [Технологический стек](#технологический-стек)
- [Архитектура](#архитектура)
- [Структура проекта](#структура-проекта)
- [Модель данных](#модель-данных)
- [API Эндпоинты](#api-эндпоинты)
- [Установка и запуск](#установка-и-запуск)
- [Работа с миграциями](#работа-с-миграциями)
- [Генерация кода из OpenAPI](#генерация-кода-из-openapi)
- [Тестирование](#тестирование)
- [Известные ограничения](#известные-ограничения)

## Описание проекта

WebPet1 — учебный/пет-проект бэкенд-сервиса для управления задачами (tasks) и пользователями (users). Каждая задача принадлежит конкретному пользователю (связь один-ко-многим). API позволяет:

- регистрировать пользователей (с хешированием пароля через bcrypt);
- получать список пользователей и информацию по конкретному пользователю;
- получать все задачи конкретного пользователя;
- создавать, получать, обновлять и удалять задачи (soft delete);
- обновлять и удалять пользователей (soft delete).

## Технологический стек

| Компонент | Технология |
|---|---|
| Язык | Go 1.24.4 |
| HTTP-роутинг | стандартный `net/http` (`http.ServeMux`) |
| ORM | [GORM](https://gorm.io/) (`gorm.io/gorm`, драйвер `gorm.io/driver/postgres`) |
| База данных | PostgreSQL |
| Миграции | [golang-migrate](https://github.com/golang-migrate/migrate) |
| API-контракт | OpenAPI 3.0 + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) (strict-server, std-http-server) |
| Хеширование паролей | `golang.org/x/crypto/bcrypt` |
| Тестирование | `testify` (assert + mock) |

## Архитектура

Проект построен по принципу слоистой архитектуры, где каждый слой отвечает строго за свою зону:

```
HTTP-запрос
    │
    ▼
[ web/tasks, web/users ]      — слой HTTP-хендлеров (сгенерирован + ручной код)
    │  принимает уже распарсенные структуры из api.gen.go,
    │  НЕ работает с БД, НЕ парсит/пишет JSON, НЕ управляет роутингом
    ▼
[ taskService, userService ]  — слой бизнес-логики ("мозг")
    │  валидация, маппинг бизнес-модель ↔ БД-модель,
    │  ничего не знает о HTTP и о деталях БД
    ▼
[ repo.go внутри каждого сервиса ] — слой доступа к данным ("руки")
    │  прямые CRUD-запросы к БД через GORM,
    │  никакой бизнес-логики
    ▼
[ internal/db ]               — слой подключения к БД ("инструмент")
    │  только открывает соединение с PostgreSQL через GORM
    ▼
PostgreSQL
```

Ключевые архитектурные решения:

1. **Contract-first (OpenAPI-first) подход.** Источник истины — `openapi/openapi.yaml`. Из него `oapi-codegen` генерирует типы запросов/ответов и роутинг (`api.gen.go`), а разработчик реализует только бизнес-часть в `handlers.go`.
2. **Strict Server.** Хендлеры реализуют строго типизированный интерфейс (`PostTasksRequestObject` → `PostTasksResponseObject`), что исключает ручной парсинг JSON и ручную установку HTTP-кодов — это делает сгенерированный код.
3. **Разделение бизнес-модели и модели БД.** В каждом сервисе есть `orm.go` (GORM-модель, `...Struct`) и бизнес-модель в `service.go` (например, `Task`, `User`). Хендлеры и сервис всегда маппят одно в другое, не пропуская "сырые" GORM-структуры наружу.
4. **Репозиторий через интерфейс.** `TaskRepoInterface` / `UserRepoInterface` позволяют подменять реальный репозиторий (`TaskRepo`/`UserRepo`) на мок (`MockTaskRepo`/`MockUserRepo`) в юнит-тестах сервисного слоя — без поднятия реальной БД.
5. **Soft delete.** У задач и пользователей есть поле `deleted_at`; удаление помечает запись, а не удаляет физически.

## Структура проекта

```
.
├── cmd/
│   └── app/
│       └── main.go            # точка входа: сборка зависимостей (DI) и запуск HTTP-сервера
├── internal/
│   ├── db/
│   │   └── db.go               # подключение к PostgreSQL через GORM
│   ├── taskService/
│   │   ├── orm.go              # GORM-модель TaskStruct
│   │   ├── repo.go             # TaskRepo — CRUD к таблице task_structs
│   │   ├── service.go          # TaskService — бизнес-логика (валидация, маппинг)
│   │   ├── taskRepository_mock.go  # мок-репозиторий для тестов
│   │   └── taskService_test.go # юнит-тесты сервиса
│   ├── userService/
│   │   ├── orm.go              # GORM-модель UserStruct (+ связь Tasks)
│   │   ├── repo.go             # UserRepo — CRUD к таблице user_structs
│   │   ├── service.go          # UserService — бизнес-логика, хеширование пароля
│   │   ├── userRepository_mock.go
│   │   └── userService_test.go
│   └── web/
│       ├── tasks/
│       │   ├── api.gen.go      # сгенерировано oapi-codegen (не редактировать вручную)
│       │   └── handlers.go     # TaskHandler — связывает HTTP-контракт и TaskService
│       └── users/
│           ├── api.gen.go      # сгенерировано oapi-codegen (не редактировать вручную)
│           └── handlers.go     # UserHandler — связывает HTTP-контракт и UserService
├── migrations/                 # SQL-миграции (golang-migrate)
├── openapi/
│   ├── openapi.yaml             # OpenAPI 3.0 спецификация — источник истины для API
│   └── .openapi                 # конфиг генерации oapi-codegen
├── Makefile                     # команды для сборки, миграций, генерации, тестов, линта
├── go.mod / go.sum
└── README.md
```

## Модель данных

### `user_structs`

| Поле | Тип | Описание |
|---|---|---|
| `id` | SERIAL PK | идентификатор пользователя |
| `email` | VARCHAR(255) UNIQUE NOT NULL | email пользователя |
| `password` | VARCHAR(255) NOT NULL | bcrypt-хеш пароля |
| `created_at` / `updated_at` | TIMESTAMP | служебные поля |
| `deleted_at` | TIMESTAMP NULL | признак soft delete |

### `task_structs`

| Поле | Тип | Описание |
|---|---|---|
| `id` | SERIAL PK | идентификатор задачи |
| `task` | TEXT NOT NULL | текст задачи |
| `is_done` | BOOLEAN DEFAULT FALSE | статус выполнения |
| `user_id` | INTEGER, FK → `user_structs.id` (ON DELETE CASCADE), индекс | владелец задачи |
| `created_at` / `updated_at` | TIMESTAMP | служебные поля |
| `deleted_at` | TIMESTAMP NULL | признак soft delete |

Связь: один пользователь → много задач (`user_structs 1 — N task_structs`), при удалении пользователя связанные задачи удаляются каскадно на уровне БД (`ON DELETE CASCADE`).

## API Эндпоинты

Полная спецификация — в `openapi/openapi.yaml`. Кратко:

### Tasks

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/tasks` | получить все задачи |
| `POST` | `/tasks` | создать задачу (`task`, `is_done?`, `user_id` — обязателен) |
| `PATCH` | `/tasks/{id}` | частично обновить задачу |
| `DELETE` | `/tasks/{id}` | удалить задачу (204 / 404) |

### Users

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/users` | получить всех пользователей |
| `POST` | `/users` | зарегистрировать пользователя (`email`, `password`) |
| `PATCH` | `/users/{id}` | частично обновить пользователя |
| `DELETE` | `/users/{id}` | удалить пользователя (204 / 404) |
| `GET` | `/users/{id}/tasks` | получить все задачи конкретного пользователя |

Пример запроса на создание задачи:

```bash
curl -X POST http://localhost:9092/tasks \
  -H "Content-Type: application/json" \
  -d '{"task": "Купить хлеб", "user_id": 1}'
```

Пример запроса на регистрацию пользователя:

```bash
curl -X POST http://localhost:9092/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "supersecret"}'
```

> Пароль в ответах API никогда не возвращается (см. схему `User` в OpenAPI).

## Установка и запуск

### Требования

- Go ≥ 1.24
- PostgreSQL (локально или в контейнере)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) — для применения миграций
- (опционально) [oapi-codegen CLI](https://github.com/oapi-codegen/oapi-codegen) — только если нужно перегенерировать код из OpenAPI

### 1. Клонирование и зависимости

```bash
git clone https://github.com/AntonRadchenko/pet1.git
cd pet1
go mod download
```

### 2. Настройка базы данных

По умолчанию приложение и миграции используют строку подключения, зашитую в коде:

- `internal/db/db.go`:
  `host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable`
- `Makefile` (`DB_DSN`):
  `postgres://postgres:yourpassword@localhost:5432/postgres?sslmode=disable`

Поднимите PostgreSQL с такими же параметрами, например через Docker:

```bash
docker run --name pet1-postgres \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=postgres \
  -p 5432:5432 \
  -d postgres:16
```

> Если хотите использовать другие креды/хост — поменяйте DSN в `internal/db/db.go` и `Makefile` соответственно (в проекте пока нет чтения параметров из `.env`/переменных окружения).

### 3. Применение миграций

```bash
make migrate
```

### 4. Запуск приложения

```bash
make run
# или напрямую:
go run cmd/app/main.go
```

После старта сервер слушает `http://localhost:9092`.

## Работа с миграциями

Все команды миграций описаны в `Makefile` и используют `golang-migrate`:

```bash
# создать новую миграцию (сгенерирует .up.sql и .down.sql)
make migrate-new NAME=имя_миграции

# применить все миграции
make migrate

# откатить миграции
make migrate-down
```

## Генерация кода из OpenAPI

Файлы `internal/web/tasks/api.gen.go` и `internal/web/users/api.gen.go` — сгенерированы автоматически и не должны редактироваться вручную. При изменении `openapi/openapi.yaml` перегенерируйте код:

```bash
# только tasks
make gen-tasks

# только users
make gen-users

# оба сразу
make gen
```

Требует установленного `oapi-codegen`:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Тестирование

Юнит-тесты покрывают сервисный слой (`taskService`, `userService`) с использованием мок-репозиториев (testify/mock), без обращения к реальной БД.

```bash
make test
# или
go test ./... -v
```

Линтер (golangci-lint):

```bash
make lint
```

## Известные ограничения

Этот раздел честно фиксирует текущее состояние проекта (полезно для тех, кто будет его дорабатывать):

- **Нет аутентификации/авторизации.** Пароль хешируется при регистрации, но эндпоинта логина, JWT/сессий и middleware для проверки прав пока нет — любой клиент может обращаться к любым `user_id`.
- **DSN базы данных захардкожен** в `internal/db/db.go` и `Makefile`, не читается из `.env`/переменных окружения.
- **Нет Dockerfile / docker-compose** для контейнеризации самого приложения (только пример подключения к БД в этом README).
- **`GetUser` в `userService`** не оборачивает `sql.ErrNoRows`/`gorm.ErrRecordNotFound` в понятную ошибку в самом сервисе (маппинг на 404 сделан только для части эндпоинтов, например `DeleteUsersId`, `GetUsersIdTasks`).

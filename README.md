# Flashcard Learning Backend API

REST API для платформы карточек (flashcards) на Go.

## Предметная область
Пользователь регистрируется и работает со своими коллекциями карточек.
Внутри коллекций можно создавать карточки вопрос/ответ и запускать тренировку в случайном порядке.

Основная сущность для CRUD: `collection`.

## Технологии
- Go 1.22+
- Chi (router)
- GORM
- PostgreSQL
- JWT (auth)
- go-playground/validator
- Docker Compose

## Архитектура
Handler -> Service -> Repository -> Database

## Быстрый старт
1. Запуск БД:

```bash
docker compose up -d db
```

2. Запуск приложения:

```bash
go run ./cmd
```

По умолчанию приложение слушает `localhost:8000`.

## Переменные окружения
Используется префикс `APP_`:
- `APP_SERVER_IP`
- `APP_SERVER_PORT`
- `APP_DATABASE_HOST`
- `APP_DATABASE_PORT`
- `APP_DATABASE_USER`
- `APP_DATABASE_PASSWORD`
- `APP_DATABASE_DBNAME`
- `APP_DATABASE_SSLMODE`
- `APP_DATABASE_TIMEZONE`
- `APP_JWT_SIGNINGKEY`

## Маршруты API

### Auth
- `POST /auth/signup`
- `POST /auth/login`
- `POST /auth/logout`

### Profile (требуется Bearer token)
- `GET /profile/`
- `PUT /profile/username`
- `PUT /profile/password`

### Collections (требуется Bearer token)
- `GET /collections/`
- `POST /collections/`
- `GET /collections/{collectionID}/`
- `PUT /collections/{collectionID}/`
- `PATCH /collections/{collectionID}/`
- `DELETE /collections/{collectionID}/`
- `GET /collections/{collectionID}/train`

### Cards (требуется Bearer token)
- `GET /collections/{collectionID}/cards/`
- `POST /collections/{collectionID}/cards/`
- `PUT /cards/{cardID}/`
- `PATCH /cards/{cardID}/`
- `DELETE /cards/{cardID}/`

## Примеры запросов/ответов

Регистрация:

```http
POST /auth/signup
Content-Type: application/json

{
	"username": "alice",
	"password": "alice_pass"
}
```

```json
{
	"id": 1,
	"username": "alice",
	"token": "<jwt>"
}
```

Создание коллекции:

```http
POST /collections/
Authorization: Bearer <jwt>
Content-Type: application/json

{
	"name": "English A2",
	"description": "Basic words"
}
```

```json
{
	"id": 3,
	"name": "English A2",
	"description": "Basic words",
	"createdAt": "2026-03-25T18:00:00Z"
}
```

Частичное обновление коллекции:

```http
PATCH /collections/3/
Authorization: Bearer <jwt>
Content-Type: application/json

{
	"name": "English A2 Updated"
}
```

```json
{
	"id": 3,
	"name": "English A2 Updated",
	"description": "Basic words",
	"createdAt": "2026-03-25T18:00:00Z"
}
```

## Единый формат ошибок
Все ошибки отдаются в JSON:

```json
{
	"error": "error message"
}
```

Примеры:
- 400 Bad Request: невалидный JSON, неверный формат path-параметра
- 401 Unauthorized: отсутствует/невалидный токен
- 404 Not Found: сущность не найдена
- 405 Method Not Allowed: метод не поддерживается
- 422 Unprocessable Entity: JSON валиден, но не проходит валидацию
- 500 Internal Server Error: внутренняя ошибка

## Валидация
- обязательные поля отмечены `validate:"required"`
- ID-поля проверяются как `gt=0`
- неизвестные поля в JSON отклоняются (`DisallowUnknownFields`)

## Тесты

```bash
go test ./...
```
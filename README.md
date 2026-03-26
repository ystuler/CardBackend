# Flashcard Learning Backend API

REST API для платформы обучения по карточкам.

Ссылка на отчет: [REPORT.md](REPORT.md)

Также можно воспользоваться простеньким веб-сайтом, его репо [тут](https://github.com/ystuler/CardFrontend)

## Запуск

### Docker
Использует локальную БД

1. Поднять сервисы:

```bash
docker compose up --build
```

2. API будет доступен на:

```text
http://localhost:8000
```

### Bare metal

1. Поднять только локальную БД:

```bash
make db-up
```

2. Запустить приложение:

```bash
make run
```

## Конфигурация

Базовые значения лежат в [config/config.yaml](config/config.yaml).

Поддерживается переопределение через переменные окружения:

- APP_SERVER_IP
- APP_SERVER_PORT
- APP_DATABASE_HOST
- APP_DATABASE_PORT
- APP_DATABASE_USER
- APP_DATABASE_PASSWORD
- APP_DATABASE_DBNAME
- APP_DATABASE_SSLMODE
- APP_DATABASE_TIMEZONE
- APP_JWT_SIGNINGKEY

## Swagger

Генерация:

```bash
make swagger
```

В результате обновляется [swagger.yaml](swagger.yaml) в корне проекта.

Просмотр:

1. открыть [онлайн редактор](https://editor.swagger.io/)
2. File -> Import File
3. выбрать [swagger.yaml](swagger.yaml)

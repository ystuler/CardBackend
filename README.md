# Flashcard Learning Backend API

REST API для платформы обучения по карточкам.

Ссылка на отчет: [REPORT.md](REPORT.md)

Также можно воспользоваться простеньким веб-сайтом, его репо [тут](https://github.com/ystuler/CardFrontend)

## Запуск

### Docker
По умолчанию использует параметры БД из config/config.yaml.

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

Рекомендуемый способ настройки: редактировать [config/config.yaml](config/config.yaml).

Это основной источник параметров (особенно для database.host/database.port/database.user/database.password/database.dbname).

Переопределение через переменные окружения поддерживается и полезно как временный override (например, для CI/CD или одноразового запуска):

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

Приоритет значений:

1. Переменные окружения APP_*
2. Значения из config/config.yaml
3. Значения по умолчанию в коде

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

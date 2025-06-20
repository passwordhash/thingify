# Changelog

Все заметные изменения в проекте документируются в этом файле.

<!-- ## [0.0.1] – 2025-06-08

### Added -->

## [Unreleased]

### Added

- Базовая структура проекта:
    - `cmd/main.go` – точка входа с инициализацией конфигурации, логгера и graceful shutdown.
    - `internal/app` – сборка зависимостей и запуск сервера.
    - `internal/github` – GraphQL-клиент для получения задач пользователя.
    - `internal/converter` – преобразование данных GitHub в доменные модели.
    - `internal/service/monitor` – служба короткого опроса новых задач.
    - `internal/storage/inmemory` – in-memory хранилище для задач.
    - `internal/messaging/rabbitmq` – интеграция с RabbitMQ для публикации новых задач.
    - `queries/github/user_issues.graphql` – начальный GraphQL-запрос.
- Конфигурация:
    - Загрузка YAML-конфига и переменных окружения через `cleanenv`.
    - Параметры RabbitMQ, GitHub API и интервал опроса вынесены в `configs/local.yml`.
- Логирование:
    - Используется `log/slog` с цветным выводом через `tint`.
- Сценарии запуска:
    - `Taskfile.yml` с задачами `run`, `dev`, `up-local`.
- Инфраструктура:
    - `docker-compose.yml` и `docker-compose.dev.yml` для поднятия RabbitMQ.
- Документация:
    - `README.md` с описанием проекта и быстрой настройкой.
- Добавлен базовый `.gitignore` для Go-проекта.
- Добавлен RabbitMQ продюсер

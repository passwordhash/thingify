# Thingify

Thingify — это сервис для мониторинга новых задач (issues) пользователя на GitHub с помощью GraphQL API. Приложение периодически опрашивает GitHub, сохраняет новые задачи в хранилище и логирует появление новых задач.

## Возможности

- Получение задач пользователя через GitHub GraphQL API
- Периодический опрос новых задач с заданным интервалом
- Сохранение задач в in-memory хранилище
- Логирование новых задач
- Гибкая настройка через YAML-конфиг и переменные окружения

## Архитектура

<img src="assets/architecture.png" style="width:50%;" />

<!-- ![Архитектура приложения](resources/architecture.png) -->

- **cmd/main.go** — точка входа, инициализация приложения, логгера и graceful shutdown
- **internal/app** — сборка зависимостей и запуск серверной части
- **internal/service/monitor** — сервис мониторинга новых задач
- **internal/github** — клиент для работы с GitHub API
- **internal/storage/inmemory** — простое in-memory хранилище задач
- **internal/converter** — преобразование моделей GitHub в доменные модели приложения
- **internal/config** — загрузка и валидация конфигурации

## Быстрый старт

1. Установите зависимости:

    ```sh
    go mod download
    ```

2. Создайте файл `.env` и укажите ваш GitHub Token:

    ```
    GH_TOKEN=your_github_token
    ```

3. Запустите приложение:
    ```sh
    go run ./cmd/main.go -config=./configs/local.yml
    ```
    или через Taskfile:
    ```sh
    task run
    ```

## Настройка

- Конфигурация по умолчанию находится в `configs/local.yml`
- Переменные окружения можно переопределять через `.env` файл

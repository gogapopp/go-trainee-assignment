# go-trainee-assignment
## Протокол

Описан в файле [schema.yaml](schema.yaml) и [schema.json](schema.json).

## Как запустить (с докером)

создать .env из .env.example
```bash
cp .env.example .env
```

```bash
docker-compose up --build
```

## Остановить и удалить контейнеры

```bash
docker-сompose down
```
## Документация запросов в постмане

https://documenter.getpostman.com/view/26679053/2sAYXEFJfz

## Уточнение по тестам

Для запуска интеграционных тестов надо поменять значение в .env DATABASE_HOST с db на localhost

## Проблемы с которыми столкнулся

Я использовал для хеширования паролей bcrypt, но из за этого SLI был >400ms, я заменил на обычное хеширование sha256 стало <50ms

Не успевал написать тесты, но также заметил что удобнее было бы вынести инициализацию проекта в отдельный пакет, например в package app для более удобного написания интеграционных тестов с вызовом всех зависимостей через app.Run(), как вариант


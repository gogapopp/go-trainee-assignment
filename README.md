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

## Остановить приложение

```bash
docker-compose-down
```
## Документация запросов в постмане

https://documenter.getpostman.com/view/26679053/2sAYXEFJfz

Для запуска интеграционных тестов репозитория надо поменять значение в .env DATABASE_HOST с db на localhost

## Проблемы с которыми столкнулся

Я использовал для хеширования паролей bcrypt, но из за этого SLI был >400ms, я заменил на обычное хеширование sha256 стало <50ms

Не совсем успел написать тесты, также заметил что удобнее было бы вынести инициализацию проекта в отдельный пакет например app для более удобного написания интеграционный тестов с вызовом всех зависимостей через app.Run()

<!-- Постман коллекция: [documenter.getpostman.com](https://documenter.getpostman.com/view/2612412453/2sA123123DpC)  

## Generating code from a specification

Install [oapi-codegen](https://github.com/deepmap/oapi-codegen/) and generate:

```bash
oapi-codegen -package=handler -generate="chi-server,types,spec" api.yaml > internal/handler/api.gen.go
```

oapi-gen:
	@oapi-codegen -package=handler -generate="chi-server,types,spec" api.yaml > internal/handler/api.gen.go

https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025 -->

<!-- {"level":"error","ts":1739669209.4378111,"caller":"handlers/info.go:34","msg":"internal.http-server.handlers.info.InfoHandler: %!w(*fmt.wrapError=&{internal.service.info.GetUserInfo: internal.repository.postgres.info.GetUserInfo: no rows in result set 0xc000324ee0}) -->